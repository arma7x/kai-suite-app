package google_services

import (
	"context"
	"errors"
	"strings"
	"time"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"kai-suite/utils/global"
	"kai-suite/types"
	"github.com/tidwall/buntdb"
	_ "kai-suite/utils/logger"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

var (
	clock = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	fields = "names,phoneNumbers,emailAddresses,metadata"
	updateFields = "names,phoneNumbers,emailAddresses"
)

func GetContacts(config *oauth2.Config, account types.UserInfoAndToken) []*people.Person {
	ctx := context.Background()
	client := GetAuthClient(config, account.Token)
	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Warn("Unable to create people Client: ", err)
		return nil
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

func UpdateContacts(config *oauth2.Config, account types.UserInfoAndToken, contacts map[string]people.Person) ([]*people.Person, []*people.Person) {
	ctx := context.Background()
	client := GetAuthClient(config, account.Token)
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

func Sync(config *oauth2.Config, account types.UserInfoAndToken) {
	connections := GetContacts(config, account)
	if len(connections) > 0 {
		updateList := make(map[string]*people.Person)
		syncList := make(map[string]people.Person)
		for _, cloudCursor := range connections {
			// log.Info(i, " ", cloudCursor.Metadata.Sources[0].UpdateTime, " ", cloudCursor.Names[0].DisplayName, "\n\n")
			// log.Info(i, string(b), "\n\n")
			key := strings.Replace(cloudCursor.ResourceName, "/", ":", 1)
			if err := global.CONTACTS_DB.View(func(tx *buntdb.Tx) error {
				val, err := tx.Get(account.User.Id + ":" + key)
				if err != nil {
					updateList[key] = cloudCursor
					return err
				}
				metadata := types.Metadata{}
				if metadata_s, err := tx.Get("metadata:" + account.User.Id + ":" + key); err == nil {
					// log.Info(metadata_s)
					if parseErr := json.Unmarshal([]byte(metadata_s), &metadata); parseErr != nil {
						updateList[key] = cloudCursor
						return err
					}
				} else {
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

				if hashCloud != metadata.Hash {
					if cloudCursor.Metadata.Sources[0].UpdateTime > localCursor.Metadata.Sources[0].UpdateTime {
						updateList[key] = cloudCursor
						return errors.New("outdated local data" + cloudCursor.Metadata.Sources[0].UpdateTime + " " + cloudCursor.Names[0].UnstructuredName)
					} else if cloudCursor.Metadata.Sources[0].UpdateTime < localCursor.Metadata.Sources[0].UpdateTime {
						log.Info(cloudCursor.Metadata.Sources[0].UpdateTime, " ", localCursor.Metadata.Sources[0].UpdateTime, "\n")
						syncList[cloudCursor.ResourceName] = localCursor
						return errors.New("outdated cloud data " + cloudCursor.Metadata.Sources[0].UpdateTime + " " + localCursor.Names[0].UnstructuredName)
					}
				} else {
					// localCursor.Metadata.Sources[0].UpdateTime == cloudCursor.Metadata.Sources[0].UpdateTime
					// log.Info("OK:" + account.User.Id + ":" + key, " ", localCursor.Names[0].DisplayName, "\n")
					if key == "people:c9181097719823060915" {
						log.Info(localCursor.Names[0].DisplayName)
						//localCursor.Names[0].GivenName = "Ahmad " + time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
						//localCursor.Names[0].UnstructuredName = localCursor.Names[0].GivenName + " " + localCursor.Names[0].FamilyName
						//localCursor.Metadata.Sources[0].UpdateTime = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
						//log.Info(key, " to update ", localCursor.Names[0].UnstructuredName, "\n")
						//updateList[key] = &localCursor
					}
				}
				return nil
			}); err != nil {
				log.Warn(key, " ", err)
			}
		}
		if len(updateList) > 0 {
			global.CONTACTS_DB.Update(func(tx *buntdb.Tx) error {
				for key, value := range updateList {
					tempTime := value.Metadata.Sources[0].UpdateTime
					value.Metadata.Sources[0].UpdateTime = ""
					b2, _ := value.MarshalJSON()
					hash := sha256.Sum256(b2)
					metadata := types.Metadata{}
					if metadata_s, err := tx.Get("metadata:" + account.User.Id + ":" + key); err == nil {
						log.Info(metadata_s)
						if parseErr := json.Unmarshal([]byte(metadata_s), &metadata); parseErr != nil {
							metadata.Deleted = false
						}
					} else {
						metadata.Deleted = false
					}
					metadata.Hash = hex.EncodeToString(hash[:])
					if metadata_b, err := json.Marshal(metadata); err == nil {
						tx.Set("metadata:" + account.User.Id + ":" + key, string(metadata_b[:]), nil)
					}
					value.Metadata.Sources[0].UpdateTime = tempTime
					b, _ := value.MarshalJSON()
					tx.Set(account.User.Id + ":" + key, string(b), nil)
					// log.Info(account.User.Id + ":" + key)
				}
				return nil
			})
		}
		if len(syncList) > 0 {
			log.Info("syncList start\n")
			success, _ := UpdateContacts(config, account, syncList)
			global.CONTACTS_DB.Update(func(tx *buntdb.Tx) error {
				for _, person := range success {
					key := strings.Replace(person.ResourceName, "/", ":", 1)
					tempTime := person.Metadata.Sources[0].UpdateTime
					person.Metadata.Sources[0].UpdateTime = ""
					b2, _ := person.MarshalJSON()
					hash := sha256.Sum256(b2)
					metadata := types.Metadata{}
					if metadata_s, err := tx.Get("metadata:" + account.User.Id + ":" + key); err == nil {
						log.Info(metadata_s)
						if parseErr := json.Unmarshal([]byte(metadata_s), &metadata); parseErr != nil {
							metadata.Deleted = false
						}
					} else {
						metadata.Deleted = false
					}
					metadata.Hash = hex.EncodeToString(hash[:])
					if metadata_b, err := json.Marshal(metadata); err == nil {
						tx.Set("metadata:" + account.User.Id + ":" + key, string(metadata_b[:]), nil)
					}

					person.Metadata.Sources[0].UpdateTime = tempTime
					b, _ := person.MarshalJSON()
					tx.Set(account.User.Id + ":" + key, string(b), nil)
				}
				return nil
			})
			log.Info("syncList end\n")
		}
	}
	global.CONTACTS_DB.Shrink()
}
