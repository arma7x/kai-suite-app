package navigations

import (
	"sort"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"kai-suite/utils/global"
	"kai-suite/utils/websockethub"
	"kai-suite/utils/google_services"
	"kai-suite/types"
	log "github.com/sirupsen/logrus"
	custom_widget "kai-suite/widgets"
)

var (
	accountsContainer *fyne.Container
	viewContactsList func(string, string)
)

func renderGoogleAccountCards(accountsContainer *fyne.Container, accounts map[string]*types.UserInfoAndToken) {
	log.Info("Google Services Rendered")
	accountsContainer.Objects = nil
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
			custom_widget.NewButton(namespace, "Sync Cloud", func(scope string) {
				log.Info("Sync Cloud ", accounts[scope].User.Id)
				if authConfig, err := google_services.GetConfig(); err == nil {
					if token, err := google_services.RefreshToken(google_services.TokenRepository[accounts[scope].User.Id].Token); err == nil {
						google_services.TokenRepository[accounts[scope].User.Id].Token = token
						if err := google_services.Sync(authConfig, google_services.TokenRepository[accounts[scope].User.Id], RemoveContact); err != nil {
							log.Warn(err.Error())
						}
					} else {
						log.Warn(err.Error())
					}
				}
			}),
			custom_widget.NewButton(namespace, "Sync KaiOS", func(scope string) {
				log.Info("Sync KaiOS ", scope)
				websockethub.SyncContacts(scope)
			}),
			custom_widget.NewButton(namespace, "Restore Contacts", func(scope string) {
				log.Info("Restore Contacts ", accounts[scope].User.Id)
				websockethub.RestoreContact(scope)
			}),
			custom_widget.NewButton(namespace, "Contact List", func(scope string) {
				log.Info("Contact List ", accounts[scope].User.Id)
				viewContactsList(accounts[scope].User.Email + " Contacts", scope)
			}),
			widget.NewButton("Remove", func() {
				log.Info("Remove ", accounts[namespace].User.Id)
				google_services.RemoveAccount(accounts[namespace].User.Id)
				renderGoogleAccountCards(accountsContainer, accounts)
			}),
		))
		accountsContainer.Add(card)
	}
	accountsContainer.Refresh()
}

func RenderGoogleAccountContent(c *fyne.Container, viewContactsListCb func(string, string)) {
	viewContactsList = viewContactsListCb
	c.Objects = nil
	accountsContainer = container.NewAdaptiveGrid(3)
	renderGoogleAccountCards(accountsContainer, google_services.TokenRepository)
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
								renderGoogleAccountCards(accountsContainer, google_services.TokenRepository)
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
		),
		nil, nil, nil,
		container.NewVScroll(container.NewVBox(accountsContainer)),
	)
	c.Add(box)
}
