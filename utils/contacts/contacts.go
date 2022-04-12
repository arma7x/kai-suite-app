package contacts

import (
	"io"
	"sort"
	//"time"
	"strings"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"kai-suite/types"
	"kai-suite/utils/global"
	"github.com/tidwall/buntdb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/people/v1"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/signal-golang/go-vcard"
	custom_widget "kai-suite/widgets"
)

func SortContacts(persons []*people.Person) {
	sort.Slice(persons, func(i, j int) bool {
		return persons[i].Names[0].UnstructuredName < persons[j].Names[0].UnstructuredName
	})
}

func GetContacts(namespace, filter string) []*people.Person {
	indexName := strings.Join([]string{"people", namespace}, "_")
	var persons []*people.Person
	if err := global.CONTACTS_DB.View(func(tx *buntdb.Tx) error {
		tx.Ascend(indexName, func(key, val string) bool {
			var person people.Person
			if err := json.Unmarshal([]byte(val), &person); err != nil {
				return false
			}
			if filter != "" {
				filter = strings.ToLower(filter)
				if len(person.Names) > 0 {
						if strings.Contains(strings.ToLower(person.Names[0].UnstructuredName), filter) {
						persons = append(persons, &person)
						return true
					}
				}
				if len(person.PhoneNumbers) > 0 {
					if strings.Contains(strings.ToLower(person.PhoneNumbers[0].Value), filter) {
						persons = append(persons, &person)
						return true
					}
				}
				if len(person.EmailAddresses) > 0 {
					if strings.Contains(strings.ToLower(person.EmailAddresses[0].Value), filter) {
						persons = append(persons, &person)
						return true
					}
				}
			} else {
				persons = append(persons, &person)
			}
			return true
		})
		return nil
	}); err != nil {
		log.Warn(err.Error())
	}
	return persons
}

func ImportContacts() {
	log.Info("ImportContacts")
	d := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
		if err == nil && f != nil {
			log.Info("Start Import")
			progressDialog := custom_widget.NewProgressInfinite("Synchronizing", "Please wait...", global.WINDOW)
			defer f.Close()
			if err != nil {
				log.Info("Error Import")
				progressDialog.Hide()
				dialog.ShowError(err, global.WINDOW)
				log.Warn(err)
				return
			}

			dec := vcard.NewDecoder(f)
			for {
				card, err := dec.Decode()
				if err == io.EOF {
					break
				} else if err != nil {
					log.Warn(err)
					continue
				}
				personID := global.RandomID()
				person := people.Person{}
				name := &people.Name{}
				person.Names = make([]*people.Name, 1)
				phoneNumber := &people.PhoneNumber{}
				person.PhoneNumbers = make([]*people.PhoneNumber, 1)
				emailAddress := &people.EmailAddress{}
				person.EmailAddresses = make([]*people.EmailAddress, 1)
				if card.PreferredValue(vcard.FieldFormattedName) != "" {
					name.UnstructuredName = card.PreferredValue(vcard.FieldFormattedName)
				}
				if len(card.Names()) > 0 {
					name.GivenName = card.Names()[0].GivenName
					name.FamilyName = card.Names()[0].FamilyName
				} else {
					if card.PreferredValue(vcard.FieldTelephone) != "" {
						name.GivenName = card.PreferredValue(vcard.FieldTelephone)
					} else {
						log.Warn("No contact name or phone number")
						continue
					}
				}
				if card.PreferredValue(vcard.FieldTelephone) != "" {
					phoneNumber.Type = "mobile"
					phoneNumber.Value = card.PreferredValue(vcard.FieldTelephone)
				}
				if card.PreferredValue(vcard.FieldEmail) != "" {
					emailAddress.Type = "personal"
					emailAddress.Value = card.PreferredValue(vcard.FieldEmail)
				}
				person.Names[0] = name
				person.PhoneNumbers[0] = phoneNumber
				person.EmailAddresses[0] = emailAddress
				person.ResourceName = "people/" + personID
				b, _ := person.MarshalJSON()
				hash := sha256.Sum256(b)
				metadata := types.Metadata{}
				//metadata.SyncID = personID
				//metadata.SyncUpdated = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
				metadata.Hash = hex.EncodeToString(hash[:])
				metadata.Deleted = false
				if metadata_b, err := json.Marshal(metadata); err == nil {
					// log.Info(personID)
					global.CONTACTS_DB.Update(func(tx *buntdb.Tx) error {
						key := "local:people:" + personID
						metadataKey := "metadata:local:people:" + personID
						tx.Set(key, string(b), nil)
						tx.Set(metadataKey, string(metadata_b), nil)
						return nil
					})
				}
			}
			progressDialog.Hide()
			log.Info("Finish Import")
		}
	}, global.WINDOW);
	d.Show()
}
