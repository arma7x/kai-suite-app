package main

import (
	"strings"
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
	"kai-suite/utils/websockethub"
	"kai-suite/utils/google_services"
	"kai-suite/types"
	"kai-suite/theme"
	"kai-suite/navigations"
	log "github.com/sirupsen/logrus"
	custom_widget "kai-suite/widgets"
	"kai-suite/utils/contacts"
	"github.com/tidwall/buntdb"
	"encoding/json"
	"github.com/gorilla/websocket"
)

var _ fyne.Theme = (*custom_theme.LightMode)(nil)
var _ fyne.Theme = (*custom_theme.DarkMode)(nil)
var port string = "4444"

var connectionContent *fyne.Container
var messagesContent *fyne.Container
var contactsContent *fyne.Container
var eventsContent *fyne.Container
var googleServicesContent *fyne.Container

var websocketBtnTxtChan = make(chan string)
var websocketClientVisibilityChan = make(chan bool)

var contentTitle binding.String
var deviceLabel = widget.NewLabel("Device: -")
var ipPortLabel = widget.NewLabel("Ip Address: " + getLocalIP() + ":" + port)
var buttonConnect = widget.NewButton("Connect", func() {
	if websockethub.Status == false {
		go websockethub.Start(onStatusChange)
	} else {
		websockethub.Stop(onStatusChange)
	}
})

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && strings.HasPrefix(ipnet.String(), "192.") {
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
	if websockethub.Status == false {
		contentTitle.Set("Disconnected")
	} else {
		contentTitle.Set("Connected")
	}
	c.Add(
		container.NewVBox(
			deviceLabel,
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
	if _, exist := google_services.TokenRepository[namespace]; exist == false  && namespace != "local" {
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

func genGoogleAccountCards(c *fyne.Container, accountList *fyne.Container, accounts map[string]types.UserInfoAndToken) {
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
			custom_widget.NewButton(namespace, "Sync Contact", func(name_space string) {
				log.Info("Sync Contact ", accounts[name_space].User.Id)
				if authConfig, err := google_services.GetConfig(); err == nil {
					google_services.Sync(authConfig, google_services.TokenRepository[accounts[name_space].User.Id]);
				}
			}),
			custom_widget.NewButton(namespace, "Sync Calendar", func(name_space string) {
				log.Info("Sync Calendar ", accounts[name_space].User.Id)
			}),
			custom_widget.NewButton(namespace, "Sync KaiOS Contacts", func(name_space string) {
				log.Info("Sync KaiOS Contacts ", name_space)
				if websockethub.Status == false  || websockethub.Client == nil {
					return
				}
				websockethub.ContactsSyncQueue = nil
				peoples := contacts.GetPeopleContacts(name_space)
				contacts.SortContacts(peoples)
				global.CONTACTS_DB.View(func(tx *buntdb.Tx) error {
					for _, p := range peoples {
						key := strings.Join([]string{name_space, strings.Replace(p.ResourceName, "/", ":", 1)}, ":")
						metadata := types.Metadata{}
						if metadata_s, err := tx.Get("metadata:" + key); err == nil {
							if parseErr := json.Unmarshal([]byte(metadata_s), &metadata); parseErr != nil {
								// log.Warn(idx, "-", err.Error())
								return nil
							}
						} else {
							// log.Warn(idx, "~", err.Error())
							return nil
						}
						// log.Info(idx, " success")
						websockethub.EnqueueContactSync(types.TxSyncContact{Namespace: key, Metadata: metadata, Person: p})
					}
					return nil
				})
				log.Info("Total queue: ", len(websockethub.ContactsSyncQueue))
				if item, err := websockethub.DequeueContactSync(); err == nil && websockethub.Client != nil {
					bd, _ := json.Marshal(item)
					btx, _ := json.Marshal(types.WebsocketMessageFlag {Flag: 1, Data: string(bd)})
					if err := websockethub.Client.GetConn().WriteMessage(websocket.TextMessage, btx); err != nil {
						log.Error("write:", err)
					}
				}
			}),
			custom_widget.NewButton(namespace, "Sync KaiOS Calendar", func(name_space string) {
				log.Info("Sync KaiOS Calendar ", accounts[name_space].User.Id)
			}),
			custom_widget.NewButton(namespace, "Contact List", func(name_space string) {
				log.Info("Contact List ", accounts[name_space].User.Id)
				renderContactsList(accounts[name_space].User.Email + " Contacts", name_space)
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
	addr, err := global.CheckIPAddress(getLocalIP(), port);
	if err != nil {
		log.Warn(err.Error())
	} else {
		log.Info(addr)
		ipPortLabel.SetText("Ip Address: " + addr)
		websockethub.Init(addr, websocketClientVisibilityChan)
	}
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

