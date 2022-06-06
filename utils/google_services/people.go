package google_services

import (
	"context"
	"errors"
	"strings"
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
	fields = "names,phoneNumbers,emailAddresses,metadata"
	updateFields = "names,phoneNumbers,emailAddresses"
)

func GetContacts(config *oauth2.Config, account *types.UserInfoAndToken) ([]*people.Person, error) {
	ctx := context.Background()
	client := GetAuthClient(config, account.Token)
	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Warn("Unable to create people Client: ", err)
		return nil, err
	}

	run := true;
	var connections []*people.Person // type Person struct
	var r *people.ListConnectionsResponse
	var loopError error
	r, loopError = srv.People.Connections.List("people/me").PageSize(1000).PersonFields(fields).Do()
	for (run) {
		if loopError != nil {
			log.Warn("Unable to retrieve people: ", loopError)
			run = false
		} else {
			if r.NextPageToken != "" {
				log.Info(r.NextPageToken)
			}
			connections = append(connections, r.Connections...)
			if r.NextPageToken == "" {
				run = false
			} else {
				r, loopError = srv.People.Connections.List("people/me").PageSize(20).PersonFields(fields).PageToken(r.NextPageToken).Do()
			}
		}
	}
	return connections, loopError
}

func UpdateContacts(config *oauth2.Config, account *types.UserInfoAndToken, contacts map[string]*people.Person) (success []*people.Person, fail []*people.Person) {
	ctx := context.Background()
	client := GetAuthClient(config, account.Token)
	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Warn("Unable to create people Client: ", err)
		return
	}
	for key, pr := range contacts {
		if p, err := srv.People.UpdateContact(key, pr).PersonFields(fields).UpdatePersonFields(updateFields).Do(); err == nil {
			success = append(success, p)
		} else {
			fail = append(fail, p)
		}
	}
	return
}

func DeleteContacts(config *oauth2.Config, account *types.UserInfoAndToken, contacts map[string]*people.Person) (success []*people.Person, fail []*people.Person) {
	ctx := context.Background()
	client := GetAuthClient(config, account.Token)
	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Warn("Unable to create people Client: ", err)
		return
	}
	for _, pr := range contacts {
		if _, err := srv.People.DeleteContact(pr.ResourceName).Do(); err == nil {
			success = append(success, pr)
		} else {
			fail = append(fail, pr)
		}
	}
	return
}

func SyncPeople(config *oauth2.Config, account *types.UserInfoAndToken, removeContactCb func(string, *people.Person)) error {
	if connections, err := GetContacts(config, account); err != nil {
		return err
	} else if len(connections) > 0 {
		personList := make(map[string]*people.Person)
		deleteList := make(map[string]*people.Person)
		updateList := make(map[string]*people.Person)
		syncList := make(map[string]*people.Person)
		for _, cloudCursor := range connections {
			key := strings.Replace(cloudCursor.ResourceName, "/", ":", 1)
			personList[key] = cloudCursor
			if err := global.CONTACTS_DB.View(func(tx *buntdb.Tx) error {

				val, hasErr := tx.Get(account.User.Id + ":" + key)

				metadata := types.Metadata{}
				if metadata_s, err := tx.Get("metadata:" + account.User.Id + ":" + key); err == nil {
					if err := json.Unmarshal([]byte(metadata_s), &metadata); err != nil {
						updateList[key] = cloudCursor
						return err
					}
				} else {
					updateList[key] = cloudCursor
					return err
				}
				if metadata.Deleted == true && hasErr != nil {
					deleteList[key] = cloudCursor
					return errors.New("deleted local data" + cloudCursor.Metadata.Sources[0].UpdateTime + " " + cloudCursor.Names[0].UnstructuredName)
				} else if metadata.Deleted == true && hasErr == nil {
					updateList[key] = cloudCursor
					return errors.New("restore local data" + cloudCursor.Metadata.Sources[0].UpdateTime + " " + cloudCursor.Names[0].UnstructuredName)
				}

				if hasErr != nil {
					updateList[key] = cloudCursor
					return hasErr
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
						syncList[cloudCursor.ResourceName] = &localCursor
						return errors.New("outdated cloud data " + cloudCursor.Metadata.Sources[0].UpdateTime + " " + localCursor.Names[0].UnstructuredName)
					}
				}
				return nil
			}); err != nil {
				log.Warn(key, " ", err)
			}
		}
		log.Info("updateList: ", len(updateList))
		if len(updateList) > 0 {
			log.Info("updateList start")
			global.CONTACTS_DB.Update(func(tx *buntdb.Tx) error {
				for key, value := range updateList {
					tempTime := value.Metadata.Sources[0].UpdateTime
					value.Metadata.Sources[0].UpdateTime = ""
					b2, _ := value.MarshalJSON()
					hash := sha256.Sum256(b2)
					metadata := types.Metadata{}
					if metadata_s, err := tx.Get("metadata:" + account.User.Id + ":" + key); err == nil {
						if err := json.Unmarshal([]byte(metadata_s), &metadata); err == nil {
							//if metadata.Deleted == true {
							//	deleteList[key] = value
							//	continue
							//}
							metadata.Deleted = false
						} else {
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
				}
				return nil
			})
			log.Info("updateList end")
		}
		log.Info("syncList: ", len(syncList))
		if len(syncList) > 0 {
			log.Info("syncList start")
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
						if err := json.Unmarshal([]byte(metadata_s), &metadata); err != nil {
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
			log.Info("syncList end")
		}
		log.Info("deleteList: ", len(deleteList))
		if len(deleteList) > 0 {
			log.Info("deleteList start")
			success, _ := DeleteContacts(config, account, deleteList)
			global.CONTACTS_DB.Update(func(tx *buntdb.Tx) error {
				for _, person := range success {
					key := strings.Replace(person.ResourceName, "/", ":", 1)
					tx.Delete("metadata:" + account.User.Id + ":" + key)
					tx.Delete(account.User.Id + ":" + key)
					removeContactCb(account.User.Id, person)
				}
				return nil
			})
			log.Info("deleteList end")
		}
		// log.Info(len(personList))
		metadataIndexName := strings.Join([]string{"metadata", account.User.Id}, "_")
		global.CONTACTS_DB.Update(func(tx *buntdb.Tx) error {
			keys := make(map[string]string)
			tx.Ascend(metadataIndexName, func(key, value string) bool {
				find := strings.Replace(key, "metadata:" + account.User.Id + ":", "", 1)
				if _, exist := personList[find]; exist == false {
					keys[key] = value
				}
				return true
			})
			log.Info("softDelete: ", len(keys))
			if len(keys) > 0 {
				for key, value := range keys {
					metadata := types.Metadata{}
					if err := json.Unmarshal([]byte(value), &metadata); err == nil {
						metadata.Deleted = true
					}
					if metadata_b, err := json.Marshal(metadata); err == nil {
						if _, _, err := tx.Set(key, string(metadata_b[:]), nil); err != nil {
							log.Warn(err.Error())
						}
					}
				}
			}
			return nil
		})
		global.CONTACTS_DB.Shrink()
		return nil
	}
	return errors.New("Unknown Error")
}
