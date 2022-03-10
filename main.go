package main

import (
	"sort"
	"net"
	"strconv"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/data/binding"
	"kai-suite/utils/global"
	_ "kai-suite/utils/logger"
	"kai-suite/utils/websocketserver"
	"kai-suite/utils/google_services"
	"kai-suite/types/misc"
	"kai-suite/theme"
	"kai-suite/navigations"
	log "github.com/sirupsen/logrus"
	custom_widget "kai-suite/widgets"
)

var _ fyne.Theme = (*custom_theme.LightMode)(nil)
var _ fyne.Theme = (*custom_theme.DarkMode)(nil)
var port string = "4444"

var connectionContent *fyne.Container
var messagesContent *fyne.Container
var contactsContent *fyne.Container
var eventsContent *fyne.Container
var googleServicesContent *fyne.Container

type ContactCardCache struct {
	Hash string
	Card fyne.CanvasObject
}

var contentTitle binding.String
var buttonText = make(chan string)
var ipPortLabel = widget.NewLabel("Ip Address: " + getLocalIP() + ":" + port)
var buttonConnect = widget.NewButton("Connect", func() {
	if websocketserver.Status == false {
		addr, err := global.CheckIPAddress(getLocalIP(), port);
		if err != nil {
			log.Warn(err.Error())
			return
		}
		log.Info(addr)
		ipPortLabel.SetText("Ip Address: " + addr)
		go websocketserver.Start(addr, onStatusChange)
	} else {
		websocketserver.Stop(onStatusChange)
	}
})

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func onStatusChange(status bool, err error) {
	if (status) {
		log.Info("Connected")
		contentTitle.Set("Connected")
		buttonText <- "Disconnect"
	} else {
		contentTitle.Set("Disconnected")
		log.Info("Disconnected")
		buttonText <- "Connect"
	}
	if err != nil {
		log.Warn(strconv.FormatBool(status), err.Error())
	}
}

func genDummyCards(list *fyne.Container) {
	list.Objects = nil
	for i := 1; i <= 10; i++ {
		card := &widget.Card{}
		card.SetTitle("Title")
		card.SetSubTitle("Subtitle")
		list.Add(card)
	}
}

func renderConnectContent(c *fyne.Container) {
	connectionContent.Show()
	messagesContent.Hide()
	contactsContent.Hide()
	eventsContent.Hide()
	googleServicesContent.Hide()
	c.Objects = nil
	if websocketserver.Status == false {
		contentTitle.Set("Disconnected")
	} else {
		contentTitle.Set("Connected")
	}
	c.Add(
		container.NewVBox(
			ipPortLabel,
			buttonConnect,
		),
	)
}

func renderMessagesContent(c *fyne.Container) {
	connectionContent.Hide()
	messagesContent.Show()
	contactsContent.Hide()
	eventsContent.Hide()
	googleServicesContent.Hide()
	c.Objects = nil
	contentTitle.Set("Messages")
	c.Add(
		container.NewVBox(
			widget.NewLabel("Messages Content"),
		),
	)
}

func renderContactsList(title, namespace string) {
	if _, exist := google_services.TokenRepository[namespace]; exist == false {
		return
	}
	connectionContent.Hide()
	messagesContent.Hide()
	contactsContent.Show()
	eventsContent.Hide()
	googleServicesContent.Hide()
	contentTitle.Set(title)
	personsArr := contacts.GetPeopleContacts(namespace)
	navigations.RenderContactsList(namespace, personsArr)
}

func renderCalendarsContent(c *fyne.Container) {
	connectionContent.Hide()
	messagesContent.Hide()
	contactsContent.Hide()
	eventsContent.Show()
	googleServicesContent.Hide()
	c.Objects = nil
	contentTitle.Set("Calendars")
	c.Add(
		container.NewVBox(
			widget.NewLabel("Calendars Content"),
		),
	)
}

func genGoogleAccountCards(c *fyne.Container, accountList *fyne.Container, accounts map[string]misc.UserInfoAndToken) {
	accountList.Objects = nil
	namespaceArr := make([]string, 0, len(accounts))
	for name := range accounts {
		namespaceArr = append(namespaceArr, name)
	}
	sort.Strings(namespaceArr)
	for _, namespace := range namespaceArr {
		card := &widget.Card{}
		card.SetTitle(accounts[namespace].User.Name)
		card.SetSubTitle(accounts[namespace].User.Email)
		card.SetContent(container.NewAdaptiveGrid(
			2,
			custom_widget.NewButton(namespace, "Sync Contact", func(idx string) {
				log.Info("Sync Contact ", accounts[idx].User.Id)
				if authConfig, err := google_services.GetConfig(); err == nil {
					google_services.Sync(authConfig, google_services.TokenRepository[accounts[idx].User.Id]);
				}
			}),
			widget.NewButton("Sync Calendar", func() {
				log.Info("Sync Calendar ", accounts[namespace].User.Id)
			}),
			custom_widget.NewButton(namespace, "Contact List", func(idx string) {
				log.Info("Contact List ", accounts[idx].User.Id)
				renderContactsList(accounts[idx].User.Email + " Contacts", idx)
			}),
			widget.NewButton("Calendar Events", func() {
				log.Info("Calendar Events ", accounts[namespace].User.Id)
			}),
			widget.NewButton("Remove", func() {
				log.Info("Remove ", accounts[namespace].User.Id)
			}),
			widget.NewButton("Remove(all data)", func() {
				log.Info("Remove(all data) ", accounts[namespace].User.Id)
			}),
		))
		accountList.Add(card)
	}
}

func renderGAContent(c *fyne.Container) {
	connectionContent.Hide()
	messagesContent.Hide()
	contactsContent.Hide()
	eventsContent.Hide()
	googleServicesContent.Show()
	c.Objects = nil
	contentTitle.Set("Google Account")
	accountList := container.NewAdaptiveGrid(3)
	genGoogleAccountCards(c, accountList, google_services.TokenRepository)
	box := container.NewBorder(
		container.NewHBox(
			widget.NewButton("Add Google Account", func() {
				if authConfig, err := google_services.GetConfig(); err == nil {
					if err := google_services.GetTokenFromWeb(authConfig); err == nil {
						var authCode string
						d := dialog.NewEntryDialog("Auth Token", "Token", func(str string) {
							authCode = str
						}, global.WINDOW)
						d.SetOnClosed(func() {
							if _, err := google_services.SaveToken(authConfig, authCode); err == nil {
								log.Info("TokenRepository: ",len(google_services.TokenRepository))
								genGoogleAccountCards(c, accountList, google_services.TokenRepository)
							} else {
								log.Warn(err)
							}
						})
						d.Show()
					} else {
						log.Warn(err)
					}
				}
			}),
			widget.NewButton("Local Contacts", func() {}),
			widget.NewButton("Local Events", func() {}),
		),
		nil, nil, nil,
		container.NewVScroll(container.NewVBox(accountList)),
	)
	c.Add(box)
}

func main() {
	contentTitle = binding.NewString()
	contentTitle.Set("")
	contentLabel := widget.NewLabelWithData(contentTitle)
	go func() {
		for {
			select {
				case txt := <- buttonText:
					buttonConnect.SetText(txt)
			}
		}
	}()
	defer global.CONTACTS_DB.Close()
	log.Info("main", global.ROOT_PATH)
	app := app.New()
	global.WINDOW = app.NewWindow("Kai Suite")
	global.WINDOW.Resize(fyne.NewSize(800, 600))
	fyne.CurrentApp().Settings().SetTheme(&custom_theme.LightMode{})
	var menuButton *fyne.Container = container.NewVBox(
		widget.NewButton("Connection", func() {
			renderConnectContent(connectionContent)
		}),
		widget.NewButton("Messages", func() {
			renderMessagesContent(messagesContent)
		}),
		widget.NewButton("Contacts", func() {
			renderContactsList("Local Contacts", "local")
		}),
		widget.NewButton("Calendars", func() {
			renderCalendarsContent(eventsContent)
		}),
		widget.NewButton("Google Account", func() {
			renderGAContent(googleServicesContent)
		}),
	)
	menuBox := container.NewVScroll(menuButton)
	menu := container.NewMax()
	menu.Add(menuBox)

	connectionContent = container.NewMax()
	messagesContent = container.NewMax()
	contactsContent = container.NewMax()
	eventsContent = container.NewMax()
	googleServicesContent = container.NewMax()

	navigations.RenderContactsContent(contactsContent)
	renderConnectContent(connectionContent)

	global.WINDOW.SetContent(container.NewBorder(
		nil,
		nil,
		container.NewBorder(widget.NewLabel("KaiOS PC Suite"), nil, nil, nil, menu),
		nil,
		container.NewBorder(contentLabel, nil, nil, nil, connectionContent, messagesContent, contactsContent, eventsContent, googleServicesContent)),
	)
	global.WINDOW.ShowAndRun()
}

