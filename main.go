package main

import (
	"net"
	"strconv"
	"strings"
	"errors"
	_ "time"
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"kai-suite/utils/global"
	_ "kai-suite/utils/logger"
	"kai-suite/utils/websocketserver"
	"kai-suite/utils/google_services"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
	"google.golang.org/api/people/v1"
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
					connections := google_services.GetContacts(google_services.AuthInstance)
					if len(connections) > 0 {
						updateList := make(map[string]string)
						syncList := make(map[string]people.Person)
						for _, cloudCursor := range connections {
							// log.Info(i, " ", cloudCursor.Metadata.Sources[0].UpdateTime, " ", cloudCursor.Names[0].DisplayName, "\n\n")
							// log.Info(i, string(b), "\n\n")
							key := strings.Replace(cloudCursor.ResourceName, "/", ":", 1)
							if err := global.DB.View(func(tx *buntdb.Tx) error {
								val, err := tx.Get(key)
								if err != nil {
									b, _ := cloudCursor.MarshalJSON()
									updateList[key] = string(b)
									return err
								}
								var localCursor people.Person
								if err := json.Unmarshal([]byte(val), &localCursor); err != nil {
									return err
								}

								if cloudCursor.Metadata.Sources[0].UpdateTime > localCursor.Metadata.Sources[0].UpdateTime {
									b, _ := cloudCursor.MarshalJSON()
									updateList[key] = string(b)
									return errors.New("outdated local data" + cloudCursor.Metadata.Sources[0].UpdateTime + " " + cloudCursor.Names[0].GivenName)
								} else if cloudCursor.Metadata.Sources[0].UpdateTime < localCursor.Metadata.Sources[0].UpdateTime {
									log.Info(cloudCursor.Metadata.Sources[0].UpdateTime, " ", localCursor.Metadata.Sources[0].UpdateTime, "\n")
									syncList[cloudCursor.ResourceName] = localCursor
									return errors.New("outdated cloud data " + cloudCursor.Metadata.Sources[0].UpdateTime + " " + cloudCursor.Names[0].GivenName)
								} else {
									log.Info(key, " ", localCursor.Metadata.Sources[0].UpdateTime == cloudCursor.Metadata.Sources[0].UpdateTime, "\n")
									//if key == "people:c9181097719823060915" {
										//localCursor.Names[0].GivenName = "Ahmad " + time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
										//localCursor.Metadata.Sources[0].UpdateTime = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
										//log.Info(key, " to update ", localCursor.Names[0].GivenName, "\n")
										//b, _ := localCursor.MarshalJSON()
										//updateList[key] = string(b)
									//}
								}
								return nil
							}); err != nil {
								log.Warn(key, " ", err)
							}
							if len(updateList) > 0 {
								global.DB.Update(func(tx *buntdb.Tx) error {
									for k, v := range updateList {
										tx.Set(k, v, nil)
									}
									return nil
								})
							}
							if len(syncList) > 0 {
								log.Info("syncList start\n")
								google_services.UpdateContacts(google_services.AuthInstance, syncList)
								log.Info("syncList end\n")
							}
						}
					}
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
	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
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

