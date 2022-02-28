package websocketserver

import (
	"fmt"
	"io"
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/gorilla/websocket"
	"context"
	"kai-suite/types/client"
	"kai-suite/types/message"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	Status bool = false
	Server http.Server
	Client *client.Client
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
	Client = client.CreateClient("", "", false, conn);
	log.Info("upgrade success")
	defer Client.GetConn().Close()
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) || err == io.EOF {
				log.Error(err)
				break
			}
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error(err)
				break
			}
			log.Error(err)
			break
		}
		switch mt {
			case websocket.TextMessage:
				wsmsg := message.ReadMessage{}
				wsmsg.UnmarshalJSON(msg);
				log.Info("recv: ", wsmsg)
		}
		//err = conn.WriteMessage(mt, msg)
		//if err != nil {
			//log.Error("write:", err)
			//break
		//}
	}
}

func Start(addr string, fn func(bool, error)) {
	Status = true
	fn(Status, nil)
	m := http.NewServeMux()
	Server = http.Server{Addr: addr, Handler: m}
	m.HandleFunc("/", handler)
	if err := Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		Status = false
		fn(Status, err)
	}
}

func Stop(fn func(bool, error)) {
	if err := Server.Shutdown(context.Background()); err != nil {
		fn(Status, err)
	} else {
		if Client != nil {
			Client.GetConn().WriteMessage(websocket.CloseMessage, []byte{})
			Client.GetConn().Close()
			Client = nil
		}
		Status = false
		fn(Status, nil)
	}
}
