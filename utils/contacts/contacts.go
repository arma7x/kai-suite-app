package contacts

import (
	"sort"
	"encoding/json"
	"kai-suite/utils/global"
	"github.com/tidwall/buntdb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/people/v1"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/layout"
)

func GetContactCards() (*fyne.Container) {
	var persons []*people.Person
	if err := global.CONTACTS_DB.View(func(tx *buntdb.Tx) error {
		tx.Ascend("key", func(key, val string) bool {
			log.Info(key, "\n")
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
	log.Info("Length > ", len(persons))
	for _, p := range persons {
		// log.Info(i, " ", p.Names[0].DisplayName)
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
		card.SetContent(container.NewHBox(
			widget.NewButton("Detail", func() {}),
			widget.NewButton("Edit", func() {}),
			widget.NewButton("Delete", func() {}),
		))
		contactCards = append(contactCards, card)
	}
	return fyne.NewContainerWithLayout(layout.NewGridLayout(4), contactCards...)
}
