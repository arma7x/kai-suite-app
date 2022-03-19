package websockethub

import(
	"errors"
	"kai-suite/types"
	log "github.com/sirupsen/logrus"
	"encoding/json"
	"github.com/gorilla/websocket"
)

var(
	ContactsSyncQueue []types.TxSyncGoogleContact
)

func EnqueueContactSync(item types.TxSyncGoogleContact, urgent bool) {
  if urgent {
		ContactsSyncQueue = append(ContactsSyncQueue, item)
	} else {
		ContactsSyncQueue = append([]types.TxSyncGoogleContact{item}, ContactsSyncQueue...)
	}
}

func DequeueContactSync() (item types.TxSyncGoogleContact, err error) {
	size := len(ContactsSyncQueue)
	if  size == 0 {
		err = errors.New("Empty")
		return
	}
	item = ContactsSyncQueue[size - 1]
	ContactsSyncQueue = ContactsSyncQueue[:size-1]
	return
}

func GetLastContactSync() (item types.TxSyncGoogleContact, err error) {
	size := len(ContactsSyncQueue)
	if size == 0 {
		err = errors.New("Empty")
		return
	}
	item = ContactsSyncQueue[size - 1]
	return 
}

func FlushContactSync() error {
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

func RestoreContact() {
  // Flag 3
}
