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
	"kai-suite/types"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
	"google.golang.org/api/people/v1"
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
	contactPage = 0
	contactMaxPage = 0
	contactPageSegment = 0
	contactPageOffset = 0
)

func RemoveContact(namespace string, person *people.Person) {
	if _, ok := contactContactCardCache[namespace][person.ResourceName]; ok {
		delete(contactContactCardCache[namespace], person.ResourceName)
	}
}

func ViewContactsList(namespace string, personsArr []*people.Person) {
	if strings.Contains(namespace, "local") {
		buttonSync.Show()
	} else {
		buttonSync.Hide()
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
							Card: contacts.MakeContactCardWidget(namespace, person),
						}
						// log.Info("NOT CACHE: ", metadata.Hash)
					} else {
						// log.Info("CACHE: ", metadata.Hash)
						if contactContactCardCache[namespace][person.ResourceName].Hash != metadata.Hash {
							contactContactCardCache[namespace][person.ResourceName].Hash = metadata.Hash
							contactContactCardCache[namespace][person.ResourceName].Card = contacts.MakeContactCardWidget(namespace, person)
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

func RenderContactsContent(c *fyne.Container, cb func(persons map[string]people.Person, metadata map[string]types.Metadata)) {
	log.Info("Contacts Rendered")
	c.Hide()
	contactContactCardCache = make(map[string]map[string]*ContactCardCache)
	contactCardsContainer = container.NewAdaptiveGrid(4)
	paginationString = binding.NewString()
	paginationString.Set("")
	buttonSync = widget.NewButton("Sync", func() {
		global.CONTACTS_DB.View(func(tx *buntdb.Tx) error {
			persons := make(map[string]people.Person)
			metadata := make(map[string]types.Metadata)
			tx.Ascend("people_local", func(key, val string) bool {
				//log.Info("person: ", key, ", ", val)
				var person people.Person
				if err := json.Unmarshal([]byte(val), &person); err != nil {
					return true
				}
				split := strings.Split(key, ":")
				//log.Info("Person Key:", split[len(split) - 1])
				persons[split[len(split) - 1]] = person //TODO
				return true
			})
			tx.Ascend("metadata_local", func(key, val string) bool {
				//log.Info("metadata_local: ", key, ", ", val)
				var mt types.Metadata
				if err := json.Unmarshal([]byte(val), &mt); err != nil {
					return true
				}
				//log.Info("Metadata Key:", mt.SyncID)
				metadata[mt.SyncID] = mt
				return true
			})
			// log.Info(len(persons), " ", len(metadata))
			cb(persons, metadata)
			return nil
		})
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
			buttonSync,
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
