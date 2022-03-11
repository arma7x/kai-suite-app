package websocketserver

import (
	"fmt"
	"io"
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/gorilla/websocket"
	"context"
	"kai-suite/types"
	"encoding/json"
)

var (
	initialized = false
	address string
	websocketClientChan chan bool
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	Status bool = false
	Server http.Server
	Client *types.Client
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Upgrade") != "websocket" && r.Header.Get("Connection") != "Upgrade" {
		fmt.Fprintf(w, "PC Suite for KaiOS device")
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade:", err)
		return
	}
	// id as time
	Client = types.CreateClient("Unknown", conn)
	websocketClientChan <- true
	log.Info("upgrade success")
	defer Client.GetConn().Close()
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) || err == io.EOF {
				log.Error(err)
				websocketClientChan <- false
				break
			}
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error(err)
				websocketClientChan <- false
				break
			}
			log.Error(err)
			websocketClientChan <- false
			break
		}
		switch mt {
			case websocket.TextMessage:
				rx := types.WebsocketMessageFlag{}
				if err := json.Unmarshal(msg, &rx); err == nil {
					switch rx.Flag {
						case 0:
							Client.SetDevice(rx.Data)
							websocketClientChan <- true
						case 2:
							data := types.RxSyncContactFlag2{}
							if err := json.Unmarshal([]byte(rx.Data), &data); err == nil {
								DequeueContactSync()
								log.Info(rx.Flag, ": ", data.Namespace, ": ", data.SyncID, ": ", data.SyncUpdated)
								if len(ContactsSyncQueue) > 0 && Client != nil {
									item, _ := GetLastContactSync()
									bd, _ := json.Marshal(item)
									btx, _ := json.Marshal(types.WebsocketMessageFlag {Flag: 1, Data: string(bd)})
									if err := Client.GetConn().WriteMessage(websocket.TextMessage, btx); err != nil {
										log.Error("write:", err)
									}
								}
							}
						case 4:
							log.Info(rx.Flag, ": ", rx.Data)
						case 6:
							log.Info(rx.Flag, ": ", rx.Data)
					}
				}
		}
		//err = conn.WriteMessage(mt, msg)
		//if err != nil {
			//log.Error("write:", err)
			//break
		//}
	}
}

func Init(addr string, clientChan chan bool) {
	initialized = true
	address = addr
	websocketClientChan = clientChan
}

func Start(fn func(bool, error)) {
	if initialized == false {
		return
	}
	Status = true
	fn(Status, nil)
	m := http.NewServeMux()
	Server = http.Server{Addr: address, Handler: m}
	m.HandleFunc("/", handler)
	if err := Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		Status = false
		fn(Status, err)
	}
}

func Stop(fn func(bool, error)) {
	if initialized == false {
		return
	}
	if err := Server.Shutdown(context.Background()); err != nil {
		fn(Status, err)
	} else {
		if Client != nil {
			Client = nil
			Client.GetConn().WriteMessage(websocket.CloseMessage, []byte{})
			Client.GetConn().Close()
		}
		Status = false
		fn(Status, nil)
	}
}
