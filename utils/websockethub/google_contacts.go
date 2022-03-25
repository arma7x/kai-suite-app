package websockethub

import(
	"strings"
	"errors"
	"kai-suite/types"
	log "github.com/sirupsen/logrus"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/tidwall/buntdb"
	"kai-suite/utils/global"
	// "google.golang.org/api/people/v1"
	"kai-suite/utils/contacts"
)

var(
	GoogleContactsQueue []types.TxSyncGoogleContact
)

func EnqueueContactSync(item types.TxSyncGoogleContact, urgent bool) {
  if urgent {
		GoogleContactsQueue = append(GoogleContactsQueue, item)
	} else {
		GoogleContactsQueue = append([]types.TxSyncGoogleContact{item}, GoogleContactsQueue...)
	}
}

func DequeueContactSync() (item types.TxSyncGoogleContact, err error) {
	size := len(GoogleContactsQueue)
	if  size == 0 {
		err = errors.New("Empty")
		return
	}
	item = GoogleContactsQueue[size - 1]
	GoogleContactsQueue = GoogleContactsQueue[:size-1]
	return
}

func GetLastContactSync() (item types.TxSyncGoogleContact, err error) {
	size := len(GoogleContactsQueue)
	if size == 0 {
		err = errors.New("Empty")
		return
	}
	item = GoogleContactsQueue[size - 1]
	return 
}

func SyncGoogleContact() error {
	if item, err := DequeueContactSync(); err == nil && Client != nil {
		bd, _ := json.Marshal(item)
		btx, _ := json.Marshal(types.WebsocketMessageFlag {Flag: 1, Data: string(bd)})
		if err := Client.GetConn().WriteMessage(websocket.TextMessage, btx); err != nil {
			log.Warn(err.Error())
			return err
		}
	}
	return nil
}

func SyncContacts(namespace string) {
	log.Info("Sync Google Contacts ", namespace)
	if Status == false  || Client == nil {
		return
	}
	peoples := contacts.GetContacts(namespace, "")
	contacts.SortContacts(peoples)
	global.CONTACTS_DB.View(func(tx *buntdb.Tx) error {
		for _, p := range peoples {
			key := strings.Join([]string{namespace, strings.Replace(p.ResourceName, "/", ":", 1)}, ":")
			metadata := types.Metadata{}
			if metadata_s, err := tx.Get("metadata:" + key); err == nil {
				if parseErr := json.Unmarshal([]byte(metadata_s), &metadata); parseErr != nil {
					return nil
				}
			} else {
				return nil
			}
			EnqueueContactSync(types.TxSyncGoogleContact{Namespace: key, Metadata: metadata, Person: p}, false)
		}
		return nil
	})
	log.Info("Total queue: ", len(GoogleContactsQueue))
	SyncGoogleContact()
}

func RestoreGoogleContact() error {
	if item, err := DequeueContactSync(); err == nil && Client != nil {
		bd, _ := json.Marshal(item)
		btx, _ := json.Marshal(types.WebsocketMessageFlag {Flag: 3, Data: string(bd)})
		if err := Client.GetConn().WriteMessage(websocket.TextMessage, btx); err != nil {
			log.Warn(err.Error())
			return err
		}
	}
	return nil
}

func RestoreContact(namespace string) {
	log.Info("Restore Google Contacts ", namespace)
	if Status == false  || Client == nil {
		return
	}
	peoples := contacts.GetContacts(namespace, "")
	contacts.SortContacts(peoples)
	global.CONTACTS_DB.View(func(tx *buntdb.Tx) error {
		for _, p := range peoples {
			key := strings.Join([]string{namespace, strings.Replace(p.ResourceName, "/", ":", 1)}, ":")
			metadata := types.Metadata{}
			if metadata_s, err := tx.Get("metadata:" + key); err == nil {
				if parseErr := json.Unmarshal([]byte(metadata_s), &metadata); parseErr != nil {
					return nil
				}
			} else {
				return nil
			}
			EnqueueContactSync(types.TxSyncGoogleContact{Namespace: key, Metadata: metadata, Person: p}, false)
		}
		return nil
	})
	log.Info("Total queue: ", len(GoogleContactsQueue))
	RestoreGoogleContact()
}
