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

func MakeContactCardWidget(namespace string, person *people.Person) fyne.CanvasObject {
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
		card.SetSubTitle(val)
	} else {
		card.SetSubTitle("-")
	}
	id := namespace + ":" + strings.Replace(person.ResourceName, "/", ":", 1)
	card.SetContent(container.NewHBox(
		custom_widget.NewButton(id, "Detail", func(scope string) {
			log.Info("Clicked detail ", scope)
		}),
		custom_widget.NewButton(id, "Edit", func(scope string) {
			log.Info("Clicked edit ", scope)
		}),
		custom_widget.NewButton(id, "Delete", func(scope string) {
			log.Info("Clicked delete ", scope)
		}),
	))
	return card
}

func SortContacts(persons []*people.Person) {
	sort.Slice(persons, func(i, j int) bool {
		return persons[i].Names[0].DisplayName < persons[j].Names[0].DisplayName
	})
}

func GetContacts(namespace string) []*people.Person {
	indexName := strings.Join([]string{"people", namespace}, "_")
	var persons []*people.Person
	if err := global.CONTACTS_DB.View(func(tx *buntdb.Tx) error {
		tx.Ascend(indexName, func(key, val string) bool {
			var person people.Person
			if err := json.Unmarshal([]byte(val), &person); err != nil {
				return false
			}
			persons = append(persons, &person)
			return true
		})
		return nil
	}); err != nil {
		log.Warn(err.Error())
	}
	return persons
}
