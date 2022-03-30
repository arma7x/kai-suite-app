package websockethub

import(
	log "github.com/sirupsen/logrus"
	"github.com/gorilla/websocket"
	"kai-suite/types"
	"encoding/json"
)

func SyncMessages() {
	if Client != nil {
		btx, _ := json.Marshal(types.WebsocketMessageFlag {Flag: 9, Data: "sync_sms"})
		if err := Client.GetConn().WriteMessage(websocket.TextMessage, btx); err != nil {
			log.Warn(err.Error())
		}
	}
}

func SendMessage(receivers []string, message string, iccId string) {
	if Client != nil {
		item := types.TxSendSMS11{
			Receivers: receivers,
			Message: message,
			IccId: iccId,
		}
		bd, _ := json.Marshal(item)
		btx, _ := json.Marshal(types.WebsocketMessageFlag {Flag: 11, Data: string(bd)})
		if err := Client.GetConn().WriteMessage(websocket.TextMessage, btx); err != nil {
			log.Warn(err.Error())
		}
	}
}

func SyncMessagesRead(id []int) {
	if Client != nil {
		item := types.TxSyncSMSRead13{
			Id: id,
		}
		bd, _ := json.Marshal(item)
		btx, _ := json.Marshal(types.WebsocketMessageFlag {Flag: 13, Data: string(bd)})
		if err := Client.GetConn().WriteMessage(websocket.TextMessage, btx); err != nil {
			log.Warn(err.Error())
		}
	}
}

func DeleteMessages(id []int) {
  log.Info("DeleteMessages: ", id)
	if Client != nil {
		item := types.TxSyncSMSDelete15{
			Id: id,
		}
		bd, _ := json.Marshal(item)
		btx, _ := json.Marshal(types.WebsocketMessageFlag {Flag: 15, Data: string(bd)})
		if err := Client.GetConn().WriteMessage(websocket.TextMessage, btx); err != nil {
			log.Warn(err.Error())
		}
	}
}
