package contacts

import (
	"sort"
	"strings"
	"encoding/json"
	"kai-suite/utils/global"
	"github.com/tidwall/buntdb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/people/v1"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	custom_widget "kai-suite/widgets"
)

func GetPeopleContactCards(namespace string) []fyne.CanvasObject {
	indexName := strings.Join([]string{namespace, "people"}, "_")
	var persons []*people.Person
	if err := global.CONTACTS_DB.View(func(tx *buntdb.Tx) error {
		tx.Ascend(indexName, func(key, val string) bool {
			// log.Info(key, "\n")
			var person people.Person
			if err := json.Unmarshal([]byte(val), &person); err != nil {
				return false
			}
			persons = append(persons, &person)
			return true
		})
		return nil
	}); err != nil {
		log.Warn(err)
	}
	sort.Slice(persons, func(i, j int) bool {
		return persons[i].Names[0].DisplayName < persons[j].Names[0].DisplayName
	})
	var contactCards []fyne.CanvasObject
	log.Info("Contacts length > ", len(persons))
	for _, p := range persons {
		card := &widget.Card{}
		if len(p.Names) > 0 {
			card.SetTitle(p.Names[0].DisplayName)
		} else {
			card.SetTitle("-")
		}
		if len(p.PhoneNumbers) > 0 {
			val := p.PhoneNumbers[0].CanonicalForm
			if val == "" {
				val = p.PhoneNumbers[0].Value
			}
			card.SetSubTitle(val)
		} else {
			card.SetSubTitle("-")
		}
		id := namespace + ":" + strings.Replace(p.ResourceName, "/", ":", 1)
		card.SetContent(container.NewHBox(
			custom_widget.NewButton(id, "Detail", func(nmsp string) {
				log.Info("Clicked detail ", nmsp)
			}),
			custom_widget.NewButton(id, "Edit", func(nmsp string) {
				log.Info("Clicked edit ", nmsp)
			}),
			custom_widget.NewButton(id, "Delete", func(nmsp string) {
				log.Info("Clicked delete ", nmsp)
			}),
		))
		contactCards = append(contactCards, card)
	}
	return contactCards
}
