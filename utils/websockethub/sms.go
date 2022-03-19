package websockethub

import(
	log "github.com/sirupsen/logrus"
	"github.com/gorilla/websocket"
	"kai-suite/types"
	"encoding/json"
)

func SyncSMS() {
	if Client != nil {
		btx, _ := json.Marshal(types.WebsocketMessageFlag {Flag: 9, Data: "sync_sms"})
		if err := Client.GetConn().WriteMessage(websocket.TextMessage, btx); err != nil {
			log.Warn(err.Error())
		}
	}
}

func SendSMS(receivers []string, message string) {
	if Client != nil {
		item := types.TxSendSMS11{
			Receivers:	receivers,
			Message:  message,
		}
		bd, _ := json.Marshal(item)
		btx, _ := json.Marshal(types.WebsocketMessageFlag {Flag: 11, Data: string(bd)})
		if err := Client.GetConn().WriteMessage(websocket.TextMessage, btx); err != nil {
			log.Warn(err.Error())
		}
	}
}
