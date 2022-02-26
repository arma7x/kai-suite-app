package main

import (
	"log"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"kai-suite/libraries/websocketserver"
	configuration "kai-suite/libraries/configuration"
)

func main() {
	log.Print("main", configuration.RootPath)
	configuration.Config.IpAddress = "192.168.43.33"
	configuration.Config.Port = "5555"
	configuration.Config.Save()

	go websocketserver.Start()

	a := app.New()
	w := a.NewWindow("Hello")
	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewLabel("Hey, I'm static"),
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	w.ShowAndRun()
}
