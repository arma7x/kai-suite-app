package websockethub

import(
	"kai-suite/utils/google_services"
	log "github.com/sirupsen/logrus"
	"github.com/gorilla/websocket"
	"kai-suite/types"
	"encoding/json"
	"google.golang.org/api/calendar/v3"
)

func InitSyncCalendar(namespace string) {
	if Status == true && Client != nil {
		syncProgressChan <- true
		item := types.TxSyncEvents17{Namespace: namespace}
		bd, _ := json.Marshal(item)
		btx, _ := json.Marshal(types.WebsocketMessageFlag {Flag: 17, Data: string(bd)})
		if err := Client.GetConn().WriteMessage(websocket.TextMessage, btx); err != nil {
			syncProgressChan <- false
			log.Warn(err.Error())
		}
	}
}

func StartSyncEvent(namespace string, unsync_events []*calendar.Event) {
	log.Info("Sync Calendars ", namespace, len(unsync_events))
	if authConfig, err := google_services.GetConfig(); err == nil {
		if token, err := google_services.RefreshToken(google_services.TokenRepository[namespace].Token); err == nil {
			google_services.TokenRepository[namespace].Token = token
			google_services.WriteTokensToFile()
			syncProgressChan <- true
			if events, err := google_services.SyncCalendar(authConfig, google_services.TokenRepository[namespace], unsync_events); err != nil {
				syncProgressChan <- false
				log.Warn(err.Error())
			} else {
				//for _, item := range events {
					//date_start := item.Start.DateTime
					//if date_start == "" {
						//date_start = item.Start.Date
					//}
					//date_end := item.End.DateTime
					//if date_end == "" {
						//date_end = item.End.Date
					//}
					//log.Printf("%v, %v, %v, %v, %v\n", item.Id, item.Summary, item.Description, date_start, date_end)
				//}
				syncProgressChan <- false
				item := types.TxSyncEvents19{Namespace: namespace, Events: events, SyncedEvents: unsync_events}
				bd, _ := json.Marshal(item)
				btx, _ := json.Marshal(types.WebsocketMessageFlag {Flag: 19, Data: string(bd)})
				Client.GetConn().WriteMessage(websocket.TextMessage, btx)
			}
		} else {
			syncProgressChan <- false
			log.Warn(err.Error())
		}
	} else {
		syncProgressChan <- false
		log.Warn(err.Error())
	}
}
