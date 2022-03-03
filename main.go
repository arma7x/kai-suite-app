package main

import (
	"net"
	"strconv"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"kai-suite/utils/global"
	_ "kai-suite/utils/logger"
	"kai-suite/utils/websocketserver"
	"kai-suite/utils/google_services"
	log "github.com/sirupsen/logrus"
)

var port string = "4444"
var content *fyne.Container
var buttonText = make(chan string)
var statusLabel = widget.NewLabel("Disconnected")
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
		statusLabel.SetText("Connected")
		buttonText <- "Disconnect"
	} else {
		statusLabel.SetText("Disconnected")
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
	content.Add(
		container.NewVBox(
			statusLabel,
			ipPortLabel,
			buttonConnect,
		),
	)
}

func renderMessagesContent() {
	for _, l := range content.Objects {
		content.Remove(l);
	}
	content.Add(
		container.NewVBox(
			widget.NewLabel("Messages"),
		),
	)
}

func renderContactsContent() {
	for _, l := range content.Objects {
		content.Remove(l);
	}
	content.Add(
		container.NewVBox(
			widget.NewLabel("Contacts"),
		),
	)
}

func renderCalendarsContent() {
	for _, l := range content.Objects {
		content.Remove(l);
	}
	content.Add(
		container.NewVBox(
			widget.NewLabel("Calendars"),
		),
	)
}

func renderGAContent() {
	for _, l := range content.Objects {
		content.Remove(l);
	}
	content.Add(
		container.NewVBox(
			widget.NewLabel("Google Account"),
		),
	)
}

func main() {
	go func() {
		for {
			select {
				case txt := <- buttonText:
					buttonConnect.SetText(txt)
			}
		}
	}()
	log.Info("main", global.ROOT_PATH)
	app := app.New()
	global.WINDOW = app.NewWindow("Kai Suite")
	var menu *fyne.Container = container.NewVBox(
		widget.NewButton("Connection", func() {
			renderConnectContent()
			content.Refresh()
		}),
		widget.NewButton("Messages", func() {
			renderMessagesContent()
			content.Refresh()
		}),
		widget.NewButton("Contacts", func() {
			if google_services.AuthInstance != nil {
				google_services.People(google_services.AuthInstance)
			}
			renderContactsContent()
			content.Refresh()
		}),
		widget.NewButton("Calendars", func() {
			if google_services.AuthInstance != nil {
				google_services.Calendar(google_services.AuthInstance)
			}
			renderCalendarsContent()
			content.Refresh()
		}),
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
			renderGAContent()
			content.Refresh()
		}),
	)
	size := menu.Size()
	size.Width = 20
	menu.Resize(size)
	content = container.NewMax()
	renderConnectContent()
	global.WINDOW.SetContent(container.NewVBox(
		container.NewHBox(widget.NewLabel("KaiOS PC Suite")),
		container.NewHBox(
			menu,
			content,
		),
	))
	global.WINDOW.ShowAndRun()
}

