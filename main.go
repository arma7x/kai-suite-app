package main

import (
	"runtime"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/data/binding"
	"kai-suite/utils/global"
	_ "kai-suite/utils/logger"
	"kai-suite/utils/websockethub"
	"kai-suite/utils/google_services"
	"kai-suite/theme"
	"kai-suite/navigations"
	// "github.com/sirupsen/logrus"
	"kai-suite/utils/contacts"
	"github.com/getlantern/systray"
	"fyne.io/fyne/v2/dialog"
	"kai-suite/resources"
	"kai-suite/types"
)

var (
	_ fyne.Theme = (*custom_theme.LightMode)(nil)
	_ fyne.Theme = (*custom_theme.DarkMode)(nil)
	contentTitle binding.String
	guidesContent *fyne.Container
	connectionContent *fyne.Container
	messagesContent *fyne.Container
	contactsContent *fyne.Container
	googleServicesContent *fyne.Container
	deviceIndicator = widget.NewButton("No Device", func() {})
	connectionIndicator = widget.NewButtonWithIcon("Connection", theme.NewThemedResource(resources.GetResource(resources.PowerOffIcon, "PowerOffIcon")), navigateConnectContent)
)

func viewContactsList(title, namespace, filter string) {
	if _, exist := google_services.TokenRepository[namespace]; exist == false  && namespace != "local" {
		return
	}
	guidesContent.Hide()
	connectionContent.Hide()
	messagesContent.Hide()
	contactsContent.Show()
	googleServicesContent.Hide()
	contentTitle.Set(title)
	personsArr := contacts.GetContacts(namespace, filter)
	navigations.ViewContactsList(namespace, personsArr)
}

func searchContacts(repository map[string]*types.UserInfoAndToken) {
	accountsNames := []string{"Local"}
	accountsMap := make(map[string]string)
	accountsMap["Local"] = "local"
	for k, v := range repository {
		accountsMap[v.User.Email] = k
		accountsNames = append(accountsNames, v.User.Email)
	}
	var searchDialog dialog.Dialog
	keyword := widget.NewEntry()
	accounts := widget.NewSelect(accountsNames, func(selected string) {})
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Keyword", Widget: keyword},
			{Text: "Source", Widget: accounts},
		},
		SubmitText: "Search",
		OnSubmit: func() {
			if namespace, exist := accountsMap[accounts.Selected]; exist == true {
				viewContactsList("Search " + accounts.Selected + " Contacts", namespace, keyword.Text)
			}
			searchDialog.Hide()
		},
	}
	searchDialog = dialog.NewCustom("Search Contacts", "Cancel", container.NewMax(form), global.WINDOW);
	searchDialog.Show()
	sz := searchDialog.MinSize()
	sz.Width = 400
	searchDialog.Resize(sz)
}

func navigateGuideContent() {
	contentTitle.Set("Guides")
	guidesContent.Show()
	connectionContent.Hide()
	messagesContent.Hide()
	contactsContent.Hide()
	googleServicesContent.Hide()
}

func navigateConnectContent() {
	contentTitle.Set("Connection")
	guidesContent.Hide()
	connectionContent.Show()
	messagesContent.Hide()
	contactsContent.Hide()
	googleServicesContent.Hide()
}

func navigateMessagesContent() {
	contentTitle.Set("Messages")
	guidesContent.Hide()
	connectionContent.Hide()
	messagesContent.Show()
	contactsContent.Hide()
	googleServicesContent.Hide()
	websockethub.SyncMessages()
	navigations.RefreshThreads()
}

func navigateGoogleServices() {
	contentTitle.Set("Google Account")
	guidesContent.Hide()
	connectionContent.Hide()
	messagesContent.Hide()
	contactsContent.Hide()
	googleServicesContent.Show()
}

func init() {
	go func() {
		for {
			select {
				case <- navigations.DeviceStatus:
					if websockethub.Client != nil {
						deviceIndicator.SetText(websockethub.Client.GetDevice())
					} else {
						deviceIndicator.SetText("No Device")
					}
				case <- navigations.ConnectionStatus:
					if (websockethub.Status) {
						connectionIndicator.SetIcon(theme.NewThemedResource(resources.GetResource(resources.PowerOnIcon, "PowerOnIcon")))
					} else {
						connectionIndicator.SetIcon(theme.NewThemedResource(resources.GetResource(resources.PowerOffIcon, "PowerOffIcon")))
					}
			}
		}
	}()
}

func main() {
	defer global.CONTACTS_DB.Close()
	contentTitle = binding.NewString()
	contentTitle.Set("")
	contentLabel := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
	contentLabel.Bind(contentTitle)
	global.APP = app.NewWithID("Kai Suite")
	global.APP.Settings().SetTheme(&custom_theme.LightMode{})
	global.APP.SetIcon(resources.GetResource(resources.AppIcon, "AppIcon"))
	global.WINDOW = global.APP.NewWindow("Kai Suite")
	global.WINDOW.Resize(fyne.NewSize(800, 600))
	var menuButton *fyne.Container = container.NewBorder(
		container.NewVBox(
			widget.NewButton("Guides", func() {
				navigateGuideContent()
			}),
			widget.NewButton("Messages", func() {
				navigateMessagesContent()
			}),
			widget.NewButton("Search Contacts", func() {
				searchContacts(google_services.TokenRepository)
			}),
			widget.NewButton("Local Contacts", func() {
				viewContactsList("Local Contacts", "local", "")
			}),
			widget.NewButton("Google Account", func() {
				navigateGoogleServices()
			}),
		),
		widget.NewButton("Quit", func() {
			global.WINDOW.Close()
		}),
		nil,nil,
	)
	menuBox := container.NewVScroll(menuButton)
	menu := container.NewMax()
	menu.Add(menuBox)

	guidesContent = container.NewMax()
	navigations.RenderGuidesContent(guidesContent)
	connectionContent = container.NewMax()
	navigations.RenderConnectionContent(connectionContent)
	googleServicesContent = container.NewMax()
	navigations.RenderGoogleAccountContent(googleServicesContent, viewContactsList)
	contactsContent = container.NewMax()
	navigations.RenderContactsContent(contactsContent, websockethub.SyncLocalContacts, websockethub.RestoreLocalContacts, contacts.ImportContacts, viewContactsList)
	messagesContent = container.NewMax()
	navigations.RenderMessagesContent(messagesContent, websockethub.SyncMessages, websockethub.SendMessage, websockethub.SyncMessagesRead, websockethub.DeleteMessages)
	navigateGuideContent()

	global.WINDOW.SetContent(container.NewBorder(
		nil,
		nil,
		container.NewBorder(
			widget.NewLabelWithStyle("KaiOS PC Suite", fyne.TextAlignLeading, fyne.TextStyle{Bold:true}),
			nil, nil, nil,
			menu,
		),
		nil,
		container.NewBorder(
			container.NewBorder(
				nil, nil,
				contentLabel,
				container.NewHBox(
					deviceIndicator,
					connectionIndicator,
					widget.NewButtonWithIcon("", theme.NewThemedResource(resources.GetResource(resources.BulbIcon, "BulbIcon")), func() {
						if global.THEME == 0 {
							global.APP.Settings().SetTheme(&custom_theme.DarkMode{})
							global.THEME = 1
						} else {
							global.APP.Settings().SetTheme(&custom_theme.LightMode{})
							global.THEME = 0
						}
					}),
				),
			),
			nil, nil, nil,
			guidesContent, connectionContent, messagesContent, contactsContent, googleServicesContent),
		),
	)
	onExit := func() {}
	global.WINDOW.SetOnClosed(exit)
	global.WINDOW.SetCloseIntercept(func() {
		global.WINDOW.Hide()
		global.VISIBILITY = false
	})
	go systray.Run(onReady, onExit)
	global.WINDOW.CenterOnScreen()
	global.WINDOW.ShowAndRun()
}

func exit() {
	websockethub.Stop()
}

func onReady() {
	var resIcon *fyne.StaticResource
	if runtime.GOOS == "windows" {
		resIcon = resources.GetResource(resources.AppIconWin, "AppIcon")
	} else {
		resIcon = resources.GetResource(resources.AppIcon, "AppIcon")
	}
	systray.SetTemplateIcon(resIcon.StaticContent, resIcon.StaticContent)
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
