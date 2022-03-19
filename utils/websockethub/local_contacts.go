package websockethub

import(
	"strings"
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/gorilla/websocket"
	"kai-suite/types"
	"encoding/json"
	"kai-suite/utils/global"
	"github.com/tidwall/buntdb"
	"google.golang.org/api/people/v1"
)

func localContactListHandler(w http.ResponseWriter, r *http.Request) {
		global.CONTACTS_DB.View(func(tx *buntdb.Tx) error {
		persons := make(map[string]people.Person)
		metadata := make(map[string]types.Metadata)
		tx.Ascend("people_local", func(key, val string) bool {
			var person people.Person
			if err := json.Unmarshal([]byte(val), &person); err != nil {
				return true
			}
			split := strings.Split(key, ":")
			persons[split[len(split) - 1]] = person //TODO
			return true
		})
		tx.Ascend("metadata_local", func(key, val string) bool {
			var mt types.Metadata
			if err := json.Unmarshal([]byte(val), &mt); err != nil {
				return true
			}
			metadata[mt.SyncID] = mt
			return true
		})
		data := types.TxSyncLocalContact5{
			Metadata:	metadata,
			Persons:  persons,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
		return nil
	})
}

func SyncLocalContacts(persons map[string]people.Person, metadata map[string]types.Metadata) {
	// log.Info(len(persons), " ", len(metadata))
	if Client != nil {
		item := types.TxSyncLocalContact5{
			Metadata:	metadata,
			Persons:  persons,
		}
		bd, _ := json.Marshal(item)
		btx, _ := json.Marshal(types.WebsocketMessageFlag {Flag: 5, Data: string(bd)})
		if err := Client.GetConn().WriteMessage(websocket.TextMessage, btx); err != nil {
			log.Warn(err.Error())
		}
	}
}

func RestoreLocalContacts(persons map[string]people.Person, metadata map[string]types.Metadata) {
	// log.Info(len(persons), " ", len(metadata))
	if Client != nil {
		item := types.TxRestoreLocalContact7{
			Metadata:	metadata,
			Persons:  persons,
		}
		bd, _ := json.Marshal(item)
		btx, _ := json.Marshal(types.WebsocketMessageFlag {Flag: 7, Data: string(bd)})
		if err := Client.GetConn().WriteMessage(websocket.TextMessage, btx); err != nil {
			log.Warn(err.Error())
		}
	}
}
