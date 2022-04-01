package navigations

import (
	"strings"
	"strconv"
	"math"
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/data/binding"
	"kai-suite/utils/global"
	"kai-suite/utils/contacts"
	"kai-suite/types"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
	"google.golang.org/api/people/v1"
	"fyne.io/fyne/v2/theme"
	custom_widget "kai-suite/widgets"
)

type ContactCardCache struct {
	Hash string
	Card fyne.CanvasObject
}

var (
	contactContactCardCache map[string]map[string]*ContactCardCache // namespace:resourceName
	contactCardsContainer *fyne.Container
	contactCards []fyne.CanvasObject
	paginationString binding.String
	paginationLabel *widget.Label
	buttonSync *widget.Button
	buttonRestore *widget.Button
	buttonImport *widget.Button
	buttonRefresh *widget.Button
	contactPage = 0
	contactMaxPage = 0
	contactPageSegment = 0
	contactPageOffset = 0
)

func renderContactItem(phone string) *custom_widget.ContextMenuButton {
	exportMenu := fyne.NewMenuItem("Copy", func() {
		global.WINDOW.Clipboard().SetContent(phone)
	})
	deleteMenu := fyne.NewMenuItem("Send SMS", func() {
		SendNewMessage(phone, "")
	})
	menu := fyne.NewMenu("", exportMenu, deleteMenu)
	return custom_widget.NewContextMenu(theme.MoreVerticalIcon(), menu)
}

func makeContactCardWidget(namespace string, person *people.Person) fyne.CanvasObject {
	card := &widget.Card{}
	if len(person.Names) > 0 {
		card.SetTitle(person.Names[0].UnstructuredName)
	} else {
		card.SetTitle("-")
	}
	if len(person.PhoneNumbers) > 0 {
		val := person.PhoneNumbers[0].CanonicalForm
		if val == "" {
			val = person.PhoneNumbers[0].Value
		}
		card.SetContent(container.NewBorder(
			nil, nil,
			widget.NewLabel(val),
			renderContactItem(val),
		))
	} else {
		card.SetSubTitle("Not available")
	}
	return card
}

func RemoveContact(namespace string, person *people.Person) {
	if _, ok := contactContactCardCache[namespace][person.ResourceName]; ok {
		delete(contactContactCardCache[namespace], person.ResourceName)
	}
}

func ViewContactsList(namespace string, personsArr []*people.Person) {
	if strings.Contains(namespace, "local") {
		buttonSync.Show()
		buttonRestore.Show()
		buttonImport.Show()
		buttonRefresh.Show()
	} else {
		buttonSync.Hide()
		buttonRestore.Hide()
		buttonImport.Hide()
		buttonRefresh.Hide()
	}
	contactCards = nil
	if contactContactCardCache[namespace] == nil {
		contactContactCardCache[namespace] = make(map[string]*ContactCardCache)
	}
	contacts.SortContacts(personsArr)
	global.CONTACTS_DB.View(func(tx *buntdb.Tx) error {
		for _, person := range personsArr {
			key := strings.Replace(person.ResourceName, "/", ":", 1)
			metadata := types.Metadata{}
			if metadata_s, err := tx.Get("metadata:" + namespace + ":" + key); err == nil {
				if err := json.Unmarshal([]byte(metadata_s), &metadata); err == nil && metadata.Deleted == false {
					if _, ok := contactContactCardCache[namespace][person.ResourceName]; !ok {
						contactContactCardCache[namespace][person.ResourceName] = &ContactCardCache{
							Hash: metadata.Hash,
							Card: makeContactCardWidget(namespace, person),
						}
						// log.Info("NOT CACHE: ", metadata.Hash)
					} else {
						// log.Info("CACHE: ", metadata.Hash)
						if contactContactCardCache[namespace][person.ResourceName].Hash != metadata.Hash {
							contactContactCardCache[namespace][person.ResourceName].Hash = metadata.Hash
							contactContactCardCache[namespace][person.ResourceName].Card = makeContactCardWidget(namespace, person)
						}
					}
					contactCards = append(contactCards, contactContactCardCache[namespace][person.ResourceName].Card)
				}
			}
		}
		return nil
	})

	paginationString.Set("")
	contactPage = 1
	contactMaxPage = int(math.Ceil(float64(len(contactCards)) / float64(40)))
	contactPageSegment = contactPage - 1
	contactPageOffset = (contactPageSegment * 40) + 40
	if contactPageOffset >= len(contactCards) {
		contactPageOffset = len(contactCards)
	}
	contactCardsContainer.Objects = contactCards[contactPageSegment * 40:contactPageOffset]
	contactCardsContainer.Refresh()
	if len(personsArr) == 0 {
		contactPage = 0
	}
	paginationString.Set(strconv.Itoa(contactPage) + "/" + strconv.Itoa(contactMaxPage))
}

func RenderContactsContent(c *fyne.Container, syncCb func(), restoreCb func(), importCb func(), refreshCb func(string, string, string)) {
	log.Info("Contacts Rendered")
	c.Hide()
	contactContactCardCache = make(map[string]map[string]*ContactCardCache)
	contactCardsContainer = container.NewAdaptiveGrid(4)
	paginationString = binding.NewString()
	paginationString.Set("")
	buttonSync = widget.NewButton("Sync", func() {
		syncCb()
	})
	buttonRestore = widget.NewButton("Restore", func() {
		restoreCb()
	})
	buttonImport = widget.NewButton("Import VCF", func() {
		log.Info("Import")
		importCb()
	})
	buttonRefresh = widget.NewButton("Refresh", func() {
		log.Info("Refresh")
		refreshCb("Local Contacts", "local", "")
	})
	paginationLabel = widget.NewLabelWithData(paginationString)
	c.Objects = nil
	paginationString.Set("")
	contactPage = 0
	contactMaxPage = 0
	contactPageSegment = 0
	contactPageOffset = 0
	if contactPageOffset >= len(contactCards) {
		contactPageOffset = len(contactCards)
	}
	paginationString.Set(strconv.Itoa(contactPage) + "/" + strconv.Itoa(contactMaxPage))
	box := container.NewBorder(
		container.NewHBox(
			widget.NewButton("Prev Page", func() {
				if contactPage - 1 <= 0 {
					return
				}
				contactPage = contactPage - 1
				contactPageSegment = contactPage - 1
				contactPageOffset = (contactPageSegment * 40) + 40
				if contactPageOffset >= len(contactCards) {
					contactPageOffset = len(contactCards)
				}
				contactCardsContainer.Objects = nil
				contactCardsContainer.Objects = contactCards[contactPageSegment * 40:contactPageOffset]
				contactCardsContainer.Refresh()
				paginationString.Set(strconv.Itoa(contactPage) + "/" + strconv.Itoa(contactMaxPage))
			}),
			layout.NewSpacer(),
			container.NewHBox(
				paginationLabel,
				buttonImport,
				buttonRefresh,
				buttonSync,
				buttonRestore,
			),
			layout.NewSpacer(),
			widget.NewButton("Next Page", func() {
				if contactPage + 1 > contactMaxPage {
					return
				}
				contactPage = contactPage + 1
				contactPageSegment = contactPage - 1
				contactPageOffset = (contactPageSegment * 40) + 40
				if contactPageOffset >= len(contactCards) {
					contactPageOffset = len(contactCards)
				}
				contactCardsContainer.Objects = nil
				contactCardsContainer.Objects = contactCards[contactPageSegment * 40:contactPageOffset]
				contactCardsContainer.Refresh()
				paginationString.Set(strconv.Itoa(contactPage) + "/" + strconv.Itoa(contactMaxPage))
			}),
		),
		nil, nil, nil,
		container.NewVScroll(container.NewVBox(contactCardsContainer)),
	)
	c.Add(box)
}
