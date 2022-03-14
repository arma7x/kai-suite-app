package websockethub

import(
	"kai-suite/types"
	log "github.com/sirupsen/logrus"
	"encoding/json"
	"github.com/gorilla/websocket"
	"google.golang.org/api/people/v1"
)

func SyncLocalContacts(persons map[string]people.Person, metadata map[string]types.Metadata) {
	// log.Info(len(persons), " ", len(metadata))
	if Client != nil {
		item := types.TxSyncContact3{
			Metadata:	metadata,
			Persons:  persons,
		}
		bd, _ := json.Marshal(item)
		btx, _ := json.Marshal(types.WebsocketMessageFlag {Flag: 3, Data: string(bd)})
		if err := Client.GetConn().WriteMessage(websocket.TextMessage, btx); err != nil {
			log.Warn(err.Error())
		}
	}
}
