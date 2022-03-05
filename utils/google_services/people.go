package google_services

import (
	"context"
	"net/http"
	"errors"
	"strings"
	_ "time"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"kai-suite/utils/global"
	"github.com/tidwall/buntdb"
	_ "kai-suite/utils/logger"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

var (
	fields = "names,phoneNumbers,emailAddresses,addresses,birthdays,metadata"
	updateFields = "names,phoneNumbers,emailAddresses,addresses,birthdays"
)

func GetContacts(client *http.Client) []*people.Person {
	ctx := context.Background()

	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Warn("Unable to create people Client: ", err)
	}

	run := true;
	var connections []*people.Person // type Person struct
	var r *people.ListConnectionsResponse
	var rErr error
	r, rErr = srv.People.Connections.List("people/me").PageSize(1000).PersonFields(fields).Do()
	for (run) {
		if rErr != nil {
			log.Warn("Unable to retrieve people: ", err)
			run = false
		} else {
			if r.NextPageToken != "" {
				log.Info(r.NextPageToken, "\n")
			}
			connections = append(connections, r.Connections...)
			if r.NextPageToken == "" {
				run = false
			} else {
				r, rErr = srv.People.Connections.List("people/me").PageSize(20).PersonFields(fields).PageToken(r.NextPageToken).Do()
			}
		}
	}
	return connections
}

func CreateContacts() {}

func UpdateContacts(client *http.Client, contacts map[string]people.Person) ([]*people.Person, []*people.Person) {
	ctx := context.Background()
	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Warn("Unable to create people Client: ", err)
	}
	var success []*people.Person
	var fail		[]*people.Person
	for key, pr := range contacts {
		if p, err := srv.People.UpdateContact(key, &pr).PersonFields(fields).UpdatePersonFields(updateFields).Do(); err == nil {
			b, _ := p.MarshalJSON()
			log.Info("SUCCESS Person: ", string(b))
			success = append(success, p)
		} else {
			log.Info("FAIL Person: ", err)
			fail = append(fail, p)
		}
	}
	return success, fail
}

func DeleteContacts() {}

func SearchContacts() {}

func Sync(client *http.Client) {
	connections := GetContacts(client)
	if len(connections) > 0 {
		updateList := make(map[string]*people.Person)
		syncList := make(map[string]people.Person)
		for _, cloudCursor := range connections {
			// log.Info(i, " ", cloudCursor.Metadata.Sources[0].UpdateTime, " ", cloudCursor.Names[0].DisplayName, "\n\n")
			// log.Info(i, string(b), "\n\n")
			key := strings.Replace(cloudCursor.ResourceName, "/", ":", 1)
			if err := global.DB.View(func(tx *buntdb.Tx) error {
				val, err := tx.Get(key)
				localHash, errH := tx.Get("hash:" + key)
				if err != nil || errH != nil {
					updateList[key] = cloudCursor
					return err
				}

				var localCursor people.Person
				if err := json.Unmarshal([]byte(val), &localCursor); err != nil {
					return err
				}

				tempTime := cloudCursor.Metadata.Sources[0].UpdateTime
				cloudCursor.Metadata.Sources[0].UpdateTime = ""
				b2, _ := cloudCursor.MarshalJSON()
				tempHash := sha256.Sum256(b2)
				hashCloud := hex.EncodeToString(tempHash[:])
				cloudCursor.Metadata.Sources[0].UpdateTime = tempTime

				if hashCloud != localHash {
					if cloudCursor.Metadata.Sources[0].UpdateTime > localCursor.Metadata.Sources[0].UpdateTime {
						updateList[key] = cloudCursor
						return errors.New("outdated local data" + cloudCursor.Metadata.Sources[0].UpdateTime + " " + cloudCursor.Names[0].GivenName)
					} else if cloudCursor.Metadata.Sources[0].UpdateTime < localCursor.Metadata.Sources[0].UpdateTime {
						log.Info(cloudCursor.Metadata.Sources[0].UpdateTime, " ", localCursor.Metadata.Sources[0].UpdateTime, "\n")
						syncList[cloudCursor.ResourceName] = localCursor
						return errors.New("outdated cloud data " + cloudCursor.Metadata.Sources[0].UpdateTime + " " + cloudCursor.Names[0].GivenName)
					}
				} else {
					log.Info(key, " ", localCursor.Metadata.Sources[0].UpdateTime == cloudCursor.Metadata.Sources[0].UpdateTime, "\n")
					if key == "people:c9181097719823060915" {
						log.Info(localCursor.Names[0].DisplayName)
						//localCursor.Names[0].GivenName = "Ahmad " + time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
						//localCursor.Names[0].UnstructuredName = localCursor.Names[0].GivenName + " " + localCursor.Names[0].FamilyName
						//localCursor.Metadata.Sources[0].UpdateTime = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
						//log.Info(key, " to update ", localCursor.Names[0].GivenName, "\n")
						//updateList[key] = &localCursor
					}
				}
				return nil
			}); err != nil {
				log.Warn(key, " ", err)
			}
			if len(updateList) > 0 {
				global.DB.Update(func(tx *buntdb.Tx) error {
					for k, v := range updateList {
						tempTime := v.Metadata.Sources[0].UpdateTime
						v.Metadata.Sources[0].UpdateTime = ""
						b2, _ := v.MarshalJSON()
						hash := sha256.Sum256(b2)
						tx.Set("hash:" + k, hex.EncodeToString(hash[:]), nil)
						v.Metadata.Sources[0].UpdateTime = tempTime
						b, _ := v.MarshalJSON()
						tx.Set(k, string(b), nil)
					}
					return nil
				})
			}
			if len(syncList) > 0 {
				log.Info("syncList start\n")
				success, _ := UpdateContacts(client, syncList)
				global.DB.Update(func(tx *buntdb.Tx) error {
					for _, person := range success {
						key := strings.Replace(person.ResourceName, "/", ":", 1)
						tempTime := person.Metadata.Sources[0].UpdateTime
						person.Metadata.Sources[0].UpdateTime = ""
						b2, _ := person.MarshalJSON()
						hash := sha256.Sum256(b2)
						tx.Set("hash:" + key, hex.EncodeToString(hash[:]), nil)
						person.Metadata.Sources[0].UpdateTime = tempTime
						b, _ := person.MarshalJSON()
						tx.Set(key, string(b), nil)
					}
					return nil
				})
				log.Info("syncList end\n")
			}
		}
	}
}
