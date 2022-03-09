package main

import (
	"net"
	"strconv"
	"math"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/data/binding"
	"kai-suite/utils/global"
	_ "kai-suite/utils/logger"
	"kai-suite/utils/websocketserver"
	"kai-suite/utils/google_services"
	"kai-suite/utils/contacts"
	"kai-suite/types/misc"
	"kai-suite/theme"
	log "github.com/sirupsen/logrus"
)

var _ fyne.Theme = (*custom_theme.LightMode)(nil)
var _ fyne.Theme = (*custom_theme.DarkMode)(nil)
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

func genDummyCards(list *fyne.Container) {
	list.Objects = nil
	for i := 1; i <= 10; i++ {
		card := &widget.Card{}
		card.SetTitle("Title")
		card.SetSubTitle("Subtitle")
		list.Add(card)
	}
}

func renderConnectContent() {
	content.Objects = nil
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
	content.Objects = nil
	contentTitle.Set("Messages")
	content.Add(
		container.NewVBox(
			widget.NewLabel("Messages Content"),
		),
	)
}

func renderContactsContent(title, namespace string) {
	if _, exist := google_services.TokenRepository[namespace]; exist == false {
		return
	}
	content.Objects = nil
	contentTitle.Set(title)
	var cards []fyne.CanvasObject
	list := container.NewAdaptiveGrid(4)
	personsArr := contacts.GetPeopleContacts(namespace)
	contacts.SortContacts(personsArr)
	for _, person := range personsArr {
		cards = append(cards, contacts.MakeContactCardWidget(namespace, person))
	}
	str := binding.NewString()
	paginationLabel := widget.NewLabelWithData(str)
	page := 1
	max := int(math.Ceil(float64(len(cards)) / float64(40)))
	seg := page - 1
	high := (seg * 40) + 40
	if high >= len(cards) {
		high = len(cards)
	}
	list.Objects = cards[seg * 40:high]
	str.Set(strconv.Itoa(page) + "/" + strconv.Itoa(max))
	box := container.NewBorder(
		container.NewHBox(
			widget.NewButton("Prev Page", func() {
				if page - 1 <= 0 {
					return
				}
				page = page - 1
				seg = page - 1
				high = (seg * 40) + 40
				if high >= len(cards) {
					high = len(cards)
				}
				list.Objects = cards[seg * 40:high]
				list.Refresh()
				str.Set(strconv.Itoa(page) + "/" + strconv.Itoa(max))
			}),
			layout.NewSpacer(),
			paginationLabel,
			layout.NewSpacer(),
			widget.NewButton("Next Page", func() {
				if page + 1 > max {
					return
				}
				page = page + 1
				seg = page - 1
				high = (seg * 40) + 40
				if high >= len(cards) {
					high = len(cards)
				}
				list.Objects = cards[seg * 40:high]
				list.Refresh()
				str.Set(strconv.Itoa(page) + "/" + strconv.Itoa(max))
			}),
		),
		nil, nil, nil,
		container.NewVScroll(container.NewVBox(list)),
	)
	content.Add(box)
	list.Refresh()
}

func renderCalendarsContent() {
	content.Objects = nil
	contentTitle.Set("Calendars")
	content.Add(
		container.NewVBox(
			widget.NewLabel("Calendars Content"),
		),
	)
}

func genGoogleAccountCards(list *fyne.Container, accounts map[string]misc.UserInfoAndToken) {
	list.Objects = nil
	for namespace, acc := range accounts {
		card := &widget.Card{}
		card.SetTitle(acc.User.Name)
		card.SetSubTitle(acc.User.Email)
		card.SetContent(container.NewAdaptiveGrid(
			2,
			widget.NewButton("Sync Contact", func() {
				log.Info("Sync Contact ", acc.User.Id)
				if authConfig, err := google_services.GetConfig(); err == nil {
					google_services.Sync(authConfig, google_services.TokenRepository[acc.User.Id]);
				}
			}),
			widget.NewButton("Sync Calendar", func() {
				log.Info("Sync Calendar ", acc.User.Id)
			}),
			widget.NewButton("Contact List", func() {
				log.Info("Contact List ", acc.User.Id)
				renderContactsContent(acc.User.Email + " Contacts", namespace)
			}),
			widget.NewButton("Calendar Events", func() {
				log.Info("Calendar Events ", acc.User.Id)
			}),
			widget.NewButton("Remove", func() {
				log.Info("Remove ", acc.User.Id)
			}),
			widget.NewButton("Remove(all data)", func() {
				log.Info("Remove(all data) ", acc.User.Id)
			}),
		))
		list.Add(card)
	}
}

func renderGAContent() {
	content.Objects = nil
	contentTitle.Set("Google Account")
	list := container.NewAdaptiveGrid(3)
	genGoogleAccountCards(list, google_services.TokenRepository)
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
								genGoogleAccountCards(list, google_services.TokenRepository)
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
		container.NewVScroll(container.NewVBox(list)),
	)
	content.Add(box)
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
			renderConnectContent()
			content.Refresh()
		}),
		widget.NewButton("Messages", func() {
			renderMessagesContent()
			content.Refresh()
		}),
		widget.NewButton("Contacts", func() {
			renderContactsContent("Local Contacts", "local")
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

