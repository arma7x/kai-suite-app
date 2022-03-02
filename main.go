package main

import (
	"net"
	"strconv"
	"net/http"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "kai-suite/utils/logger"
	"kai-suite/utils/websocketserver"
	"kai-suite/utils/configuration"
	"kai-suite/utils/google"
	log "github.com/sirupsen/logrus"
)

var port string = "4444"
var content *fyne.Container
var buttonText = make(chan string)
var statusLabel = widget.NewLabel("Disconnected")
var ipPortLabel = widget.NewLabel("Ip Address: " + getLocalIP() + ":" + port)
var buttonConnect = widget.NewButton("Connect", func() {
	if websocketserver.Status == false {
		addr, err := configuration.CheckIPAddress(getLocalIP(), port);
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

func saveIpPort(addr, port string) error {
	if _, err := configuration.CheckIPAddress(addr, port); err != nil {
		log.Warn(err.Error())
		return err
	}
	configuration.Config.IpAddress = addr
	configuration.Config.Port = port
	if err := configuration.Config.Save(); err != nil {
		log.Warn(err.Error())
		return err
	}
	return nil
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
	var googleAccount *http.Client
	go func() {
		for {
			select {
				case txt := <- buttonText:
					buttonConnect.SetText(txt)
			}
		}
	}()
	log.Info("main", configuration.RootPath)
	app := app.New()
	window := app.NewWindow("Kai Suite")
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
			if googleAccount != nil {
				google.People(googleAccount)
			}
			renderContactsContent()
			content.Refresh()
		}),
		widget.NewButton("Calendars", func() {
			if googleAccount != nil {
				google.Calendar(googleAccount)
			}
			renderCalendarsContent()
			content.Refresh()
		}),
		widget.NewButton("Google Account", func() {
			googleAccount = google.GetAuth()
			renderGAContent()
			content.Refresh()
		}),
	)
	size := menu.Size()
	size.Width = 20
	menu.Resize(size)
	content = container.NewMax()
	renderConnectContent()
	window.SetContent(container.NewVBox(
		container.NewHBox(widget.NewLabel("KaiOS PC Suite")),
		container.NewHBox(
			menu,
			content,
		),
	))
	window.ShowAndRun()
}

