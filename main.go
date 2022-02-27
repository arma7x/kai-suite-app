package main

import (
	"strconv"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "kai-suite/utils/logger"
	"kai-suite/utils/websocketserver"
	"kai-suite/utils/configuration"
	log "github.com/sirupsen/logrus"
)

var buttonText = make(chan string)
var statusLabel = widget.NewLabel("KaiOS PC Suite")
var ipPortLabel = widget.NewLabel("Ip Address: nil")
var button = widget.NewButton("Connect", func() {
	if websocketserver.Status == false {
		addr, err := configuration.CheckIPAddress(configuration.Config.IpAddress, configuration.Config.Port);
		if err != nil {
			log.Warn(err.Error())
			return
		}
		log.Info(addr)
		ipPortLabel.SetText("Ip Address: " + addr)
		go websocketserver.Start(addr, onChange)
	} else {
		websocketserver.Stop(onChange)
	}
})

func onChange(status bool, err error) {
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

func main() {
	go func() {
		for {
			select {
				case txt := <- buttonText:
					button.SetText(txt)
			}
		}
	}()
	log.Info("main", configuration.RootPath)
	app := app.New()
	window := app.NewWindow("Kai Suite")
	window.SetContent(container.NewVBox(
		statusLabel,
		ipPortLabel,
		button,
	))
	window.ShowAndRun()
}
