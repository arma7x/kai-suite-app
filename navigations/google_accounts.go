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
	viewContactsList func(string, string, string)
)

// TODO CACHE WIDGET
func renderGoogleAccountCards(accountsContainer *fyne.Container, accounts map[string]*types.UserInfoAndToken) {
	log.Info("Google Services Rendered")
	accountsContainer.Objects = nil
	sortNamespace := make([]string, 0, len(accounts))
	for name := range accounts {
		sortNamespace = append(sortNamespace, name)
	}
	sort.Strings(sortNamespace)
	for _, namespace := range sortNamespace {
		scope := namespace
		card := &widget.Card{}
		card.SetTitle(accounts[namespace].User.Name)
		card.SetSubTitle(accounts[namespace].User.Email)
		card.SetContent(container.NewVBox(
			widget.NewLabelWithStyle("Peoples", fyne.TextAlignCenter, fyne.TextStyle{Bold:true}),
			container.NewAdaptiveGrid(
				2,
				widget.NewButton("Sync Cloud", func() {
					log.Info("Sync Cloud ", accounts[scope].User.Id)
					if authConfig, err := google_services.GetConfig(); err == nil {
						if token, err := google_services.RefreshToken(google_services.TokenRepository[accounts[scope].User.Id].Token); err == nil {
							google_services.TokenRepository[accounts[scope].User.Id].Token = token
							google_services.WriteTokensToFile()
							progress := custom_widget.NewProgressInfinite("Synchronizing", "Please wait...", global.WINDOW)
							if err := google_services.SyncPeople(authConfig, google_services.TokenRepository[accounts[scope].User.Id], RemoveContact); err != nil {
								progress.Hide()
								dialog.ShowError(err, global.WINDOW)
								log.Warn(err)
							} else {
								progress.Hide()
							}
						} else {
							dialog.ShowError(err, global.WINDOW)
							log.Warn(err)
						}
					} else {
						dialog.ShowError(err, global.WINDOW)
						log.Warn(err)
					}
				}),
				widget.NewButton("Sync KaiOS", func() {
					log.Info("Sync KaiOS ", scope)
					websockethub.SyncContacts(scope)
				}),
				widget.NewButton("Restore Contacts", func() {
					log.Info("Restore Contacts ", accounts[scope].User.Id)
					websockethub.RestoreContact(scope)
				}),
				widget.NewButton("Contact List", func() {
					log.Info("Contact List ", accounts[scope].User.Id)
					viewContactsList(accounts[scope].User.Email + " Contacts", scope, "")
				}),
			),
			widget.NewLabelWithStyle("Calendar", fyne.TextAlignCenter, fyne.TextStyle{Bold:true}),
			container.NewAdaptiveGrid(
				1,
				widget.NewButton("Sync Calendar", func() {
					log.Info("Sync Calendars ", accounts[scope].User.Id)
					if _, err := google_services.GetConfig(); err == nil {
						if token, err := google_services.RefreshToken(google_services.TokenRepository[accounts[scope].User.Id].Token); err == nil {
							google_services.TokenRepository[accounts[scope].User.Id].Token = token
							google_services.WriteTokensToFile()
							websockethub.InitSyncCalendar(accounts[scope].User.Id)
						} else {
							dialog.ShowError(err, global.WINDOW)
							log.Warn(err)
						}
					} else {
						dialog.ShowError(err, global.WINDOW)
						log.Warn(err)
					}
				}),
			),
			widget.NewButton("Remove This Account", func() {
				log.Info("Remove ", accounts[scope].User.Id)
				google_services.RemoveAccount(accounts[scope].User.Id)
				renderGoogleAccountCards(accountsContainer, accounts)
			}),
		))
		accountsContainer.Add(card)
	}
	accountsContainer.Refresh()
}

func RenderGoogleAccountContent(c *fyne.Container, viewContactsListCb func(string, string, string)) {
	viewContactsList = viewContactsListCb
	c.Objects = nil
	accountsContainer = container.NewAdaptiveGrid(3)
	renderGoogleAccountCards(accountsContainer, google_services.TokenRepository)
	var tokenDialog dialog.Dialog
	box := container.NewBorder(
		container.NewHBox(
			widget.NewButton("Add Google Account", func() {
				if authConfig, err := google_services.GetConfig(); err == nil {
					if err := google_services.GetTokenFromWeb(authConfig); err == nil {
						var authCode string
						tokenDialog = dialog.NewEntryDialog("Auth Token", "Token", func(str string) {
							authCode = str
						}, global.WINDOW)
						tokenDialog.SetOnClosed(func() {
							if _, err := google_services.SaveToken(authConfig, authCode); err == nil {
								log.Info("TokenRepository: ",len(google_services.TokenRepository))
								renderGoogleAccountCards(accountsContainer, google_services.TokenRepository)
							} else {
								dialog.ShowError(err, global.WINDOW)
								log.Warn(err)
							}
						})
						tokenDialog.Show()
					} else {
						tokenDialog.Hide()
						dialog.ShowError(err, global.WINDOW)
						log.Warn(err)
					}
				} else {
					dialog.ShowError(err, global.WINDOW)
					log.Warn(err)
				}
			}),
		),
		nil, nil, nil,
		container.NewVScroll(container.NewVBox(accountsContainer)),
	)
	c.Add(box)
}
