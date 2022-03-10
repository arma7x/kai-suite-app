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
	_ "kai-suite/utils/logger"
	"kai-suite/utils/contacts"
	"kai-suite/types/misc"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
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
	contactPage = 0
	contactMaxPage = 1
	contactPageSegment = 0
	contactPageOffset = 0
)

func RenderContactsList(namespace string, repositories map[string]misc.UserInfoAndToken) {
	if _, exist := repositories[namespace]; exist == false {
		return
	}
	contactCards = nil
	if contactContactCardCache[namespace] == nil {
		contactContactCardCache[namespace] = make(map[string]*ContactCardCache)
	}
	personsArr := contacts.GetPeopleContacts(namespace)
	contacts.SortContacts(personsArr)
	global.CONTACTS_DB.View(func(tx *buntdb.Tx) error {
		for _, person := range personsArr {
			key := strings.Replace(person.ResourceName, "/", ":", 1)
			metadata := &misc.Metadata{}
			if metadata_s, err := tx.Get("metadata:" + namespace + ":" + key); err == nil {
				if err := json.Unmarshal([]byte(metadata_s), &metadata); err == nil {
					if _, ok := contactContactCardCache[namespace][person.ResourceName]; !ok {
						contactContactCardCache[namespace][person.ResourceName] = &ContactCardCache{
							Hash: metadata.Hash,
							Card: contacts.MakeContactCardWidget(namespace, person),
						}
						// log.Info("NOT CACHE: ", metadata.Hash)
					} else {
						// log.Info("CACHE: ", metadata.Hash)
						if contactContactCardCache[namespace][person.ResourceName].Hash != metadata.Hash {
							contactContactCardCache[namespace][person.ResourceName].Hash = metadata.Hash
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
	paginationString.Set(strconv.Itoa(contactPage) + "/" + strconv.Itoa(contactMaxPage))
}

func RenderContactsContent(c *fyne.Container) {
	log.Info("Contacts Rendered")
	c.Hide()
	contactContactCardCache = make(map[string]map[string]*ContactCardCache)
	contactCardsContainer = container.NewAdaptiveGrid(4)
	paginationString = binding.NewString()
	paginationString.Set("")
	paginationLabel = widget.NewLabelWithData(paginationString)
	c.Objects = nil
	paginationString.Set("")
	contactPage = 0
	contactMaxPage = 1
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
			paginationLabel,
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