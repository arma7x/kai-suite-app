package main

import (
	"net"
	"strconv"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/data/binding"
	"kai-suite/utils/global"
	_ "kai-suite/utils/logger"
	"kai-suite/utils/websocketserver"
	"kai-suite/utils/google_services"
	"kai-suite/utils/contacts"
	log "github.com/sirupsen/logrus"
)

var port string = "4444"
var content *fyne.Container
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

func renderConnectContent() {
	for _, l := range content.Objects {
		content.Remove(l);
	}
	if websocketserver.Status == false {
		contentTitle.Set("Disconnected")
	} else {
		contentTitle.Set("Connected")
	}
	content.Add(
		container.NewVBox(
			ipPortLabel,
			buttonConnect,
		),
	)
}

func renderMessagesContent() {
	for _, l := range content.Objects {
		content.Remove(l);
	}
	contentTitle.Set("Messages")
	content.Add(
		container.NewVBox(
			widget.NewLabel("Messages Content"),
		),
	)
}

func renderContactsContent() {
	for _, l := range content.Objects {
		content.Remove(l);
	}
	contentTitle.Set("Contacts")
	box := container.NewVScroll(contacts.GetContacts())
	content.Add(box)
}

func renderCalendarsContent() {
	for _, l := range content.Objects {
		content.Remove(l);
	}
	contentTitle.Set("Calendars")
	content.Add(
		container.NewVBox(
			widget.NewLabel("Calendars Content"),
		),
	)
}

func renderGAContent() {
	for _, l := range content.Objects {
		content.Remove(l);
	}
	contentTitle.Set("Google Account")
	content.Add(
		container.NewVBox(
			widget.NewButton("Google Account", func() {
				google_services.AuthInstance = google_services.GetAuth()
				if google_services.AuthInstance == nil {
					if cfg, err := google_services.GetConfig(); err == nil {
						if err := google_services.GetTokenFromWeb(cfg); err == nil {
							var authCode string
							d := dialog.NewEntryDialog("Auth Token", "Token", func(str string) {
								authCode = str
							}, global.WINDOW)
							d.SetOnClosed(func() {
								if _, err := google_services.SaveToken(cfg, global.ResolvePath("token.json"), authCode); err == nil {
									google_services.AuthInstance = google_services.GetAuth()
								} else {
									log.Warn(err)
								}
							})
							d.Show()
						} else {
							log.Warn(err)
						}
					}
				}
			}),
			widget.NewButton("Sync Contacs", func() {
				if google_services.AuthInstance != nil {
					google_services.Sync(google_services.AuthInstance)
				}
			}),
			widget.NewButton("Sync Calendars", func() {
				if google_services.AuthInstance != nil {
					google_services.Calendar(google_services.AuthInstance)
				}
			}),
		),
	)
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
	defer global.DB.Close()
	log.Info("main", global.ROOT_PATH)
	app := app.New()
	global.WINDOW = app.NewWindow("Kai Suite")
	global.WINDOW.Resize(fyne.NewSize(800, 600))
	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	var menuButton *fyne.Container = container.NewVBox(
		widget.NewButton("Connection", func() {
			renderConnectContent()
			content.Refresh()
		}),
		widget.NewButton("Messages", func() {
			renderMessagesContent()
			content.Refresh()
		}),
		widget.NewButton("Contacts", func() {
			renderContactsContent()
			content.Refresh()
		}),
		widget.NewButton("Calendars", func() {
			renderCalendarsContent()
			content.Refresh()
		}),
		widget.NewButton("Google Account", func() {
			renderGAContent()
			content.Refresh()
		}),
	)
	menuBox := container.NewVScroll(menuButton)
	menu := container.NewMax()
	menu.Add(menuBox)
	content = container.NewMax()
	renderConnectContent()
	global.WINDOW.SetContent(container.NewBorder(
		nil,
		nil,
		container.NewBorder(widget.NewLabel("KaiOS PC Suite"), nil, nil, nil, menu),
		nil,
		container.NewBorder(contentLabel, nil, nil, nil, content)),
	)
	global.WINDOW.ShowAndRun()
}

