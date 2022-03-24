package main

import (
	"net"
	"strconv"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/data/binding"
	"kai-suite/utils/global"
	_ "kai-suite/utils/logger"
	"kai-suite/utils/websockethub"
	"kai-suite/utils/google_services"
	"kai-suite/theme"
	"kai-suite/navigations"
	log "github.com/sirupsen/logrus"
	"kai-suite/utils/contacts"
	"github.com/getlantern/systray"
	// "fyne.io/systray"
	"kai-suite/icon"
)

var _ fyne.Theme = (*custom_theme.LightMode)(nil)
var _ fyne.Theme = (*custom_theme.DarkMode)(nil)

var ip 		string = "(Select network card)"
var port 	string = "4444"

var connectionContent *fyne.Container
var messagesContent *fyne.Container
var contactsContent *fyne.Container
var googleServicesContent *fyne.Container

var websocketBtnTxtChan = make(chan string)
var websocketClientVisibilityChan = make(chan bool)

var contentTitle binding.String
var deviceLabel = widget.NewLabel("Device: -")
var ipPortLabel = widget.NewLabel("Ip Address: " + ip + ":" + port)
var buttonConnect = widget.NewButton("Connect", func() {
	if websockethub.Status == false {
		addr, err := global.CheckIPAddress(ip, port)
		if err != nil {
			log.Warn(err.Error())
			ipPortLabel.SetText(err.Error());
		} else {
			ipPortLabel.SetText("Ip Address: " + addr)
			websockethub.Init(addr, websocketClientVisibilityChan, navigations.ReloadThreads, navigations.ReloadMessages, navigations.RemoveContact, navigations.RefreshThreads)
		}
		go websockethub.Start(onStatusChange)
	} else {
		websockethub.Stop(onStatusChange)
	}
})

func getNetworkCardIPAddresses() (ipList []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipList = append(ipList, ipnet.IP.String())
			}
		}
	}
	return
}

func onStatusChange(status bool, err error) {
	if (status) {
		log.Info("Connected")
		contentTitle.Set("Connected")
		websocketBtnTxtChan <- "Disconnect"
	} else {
		contentTitle.Set("Disconnected")
		log.Info("Disconnected")
		websocketBtnTxtChan <- "Connect"
	}
	if err != nil {
		log.Warn(strconv.FormatBool(status), err.Error())
	}
}

func navigateConnectContent(c *fyne.Container) {
	connectionContent.Show()
	messagesContent.Hide()
	contactsContent.Hide()
	googleServicesContent.Hide()
	inputIp := widget.NewSelect(getNetworkCardIPAddresses(), func(selected string) {
		ip = selected
		ipPortLabel.SetText("Ip Address: " + ip + ":" + port)
	})
	inputIp.PlaceHolder = ip
	inputPort := widget.NewEntry()
	inputPort.Text = port
	inputPort.OnChanged = func(val string) {
		port = val
		ipPortLabel.SetText("Ip Address: " + ip + ":" + port)
	}
	c.Objects = nil
	if websockethub.Status == false {
		contentTitle.Set("Disconnected")
	} else {
		contentTitle.Set("Connected")
	}
	c.Add(
		container.NewVBox(
			deviceLabel,
			ipPortLabel,
			inputIp,
			inputPort,
			buttonConnect,
		),
	)
}

func navigateMessagesContent(c *fyne.Container) {
	contentTitle.Set("Messages")
	connectionContent.Hide()
	messagesContent.Show()
	contactsContent.Hide()
	googleServicesContent.Hide()
	websockethub.SyncSMS()
	navigations.RefreshThreads()
}

func viewContactsList(title, namespace string) {
	if _, exist := google_services.TokenRepository[namespace]; exist == false  && namespace != "local" {
		return
	}
	connectionContent.Hide()
	messagesContent.Hide()
	contactsContent.Show()
	googleServicesContent.Hide()
	contentTitle.Set(title)
	personsArr := contacts.GetContacts(namespace)
	navigations.ViewContactsList(namespace, personsArr)
}

func navigateGoogleServices(c *fyne.Container) {
	contentTitle.Set("Google Account")
	connectionContent.Hide()
	messagesContent.Hide()
	contactsContent.Hide()
	googleServicesContent.Show()
}

func main() {
	contentTitle = binding.NewString()
	contentTitle.Set("")
	contentLabel := widget.NewLabelWithData(contentTitle)
	go func() {
		for {
			select {
				case txt := <- websocketBtnTxtChan:
					buttonConnect.SetText(txt)
				case present := <- websocketClientVisibilityChan:
					if present {
						deviceLabel.SetText("Device: " + websockethub.Client.GetDevice())
					} else {
						deviceLabel.SetText("Device: -")
					}
			}
		}
	}()
	defer global.CONTACTS_DB.Close()
	global.APP = app.New()
	global.APP.Settings().SetTheme(&custom_theme.LightMode{})
	global.APP.SetIcon(&fyne.StaticResource{ StaticName: "Icon.png", StaticContent: icon.Data})
	global.WINDOW = global.APP.NewWindow("Kai Suite")
	global.WINDOW.Resize(fyne.NewSize(800, 600))
	var menuButton *fyne.Container = container.NewVBox(
		widget.NewButton("Connection", func() {
			navigateConnectContent(connectionContent)
		}),
		widget.NewButton("Messages", func() {
			navigateMessagesContent(messagesContent)
		}),
		widget.NewButton("Local Contacts", func() {
			viewContactsList("Local Contacts", "local")
		}),
		widget.NewButton("Google Account", func() {
			navigateGoogleServices(googleServicesContent)
		}),
	)
	menuBox := container.NewVScroll(menuButton)
	menu := container.NewMax()
	menu.Add(menuBox)

	connectionContent = container.NewMax()

	googleServicesContent = container.NewMax()
	navigations.RenderGoogleAccountContent(googleServicesContent, viewContactsList)
	contactsContent = container.NewMax()
	navigations.RenderContactsContent(contactsContent, websockethub.SyncLocalContacts, websockethub.RestoreLocalContacts, contacts.ImportContacts)
	messagesContent = container.NewMax()
	navigations.RenderMessagesContent(messagesContent, websockethub.SyncSMS, websockethub.SendSMS, websockethub.SyncSMSRead)
	navigateConnectContent(connectionContent)

	global.WINDOW.SetContent(container.NewBorder(
		nil,
		nil,
		container.NewBorder(widget.NewLabel("KaiOS PC Suite"), nil, nil, nil, menu),
		nil,
		container.NewBorder(contentLabel, nil, nil, nil, connectionContent, messagesContent, contactsContent, googleServicesContent)),
	)
	onExit := func() {}
	global.WINDOW.SetCloseIntercept(func() {
		global.WINDOW.Hide()
		global.VISIBILITY = false
		// global.APP.SendNotification(fyne.NewNotification("title", "content"))
	})
	go systray.Run(onReady, onExit)
	global.WINDOW.ShowAndRun()
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	// systray.SetTitle("Kai Suite")
	systray.SetTooltip("Kai Suite")
	mLaunch := systray.AddMenuItem("Launch", "Launch app")
	mQuit := systray.AddMenuItem("Quit", "Quit app")
	for {
		select {
			case <-mLaunch.ClickedCh:
				global.WINDOW.Show()
				global.VISIBILITY = true
			case <-mQuit.ClickedCh:
				global.WINDOW.Close()
		}
	}
}
