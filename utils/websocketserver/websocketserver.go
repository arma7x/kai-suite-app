package websocketserver

import ( 
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/gorilla/websocket"
	"context"
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
)

var client *websocket.Conn

func connect(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade:", err)
		return
	}
	client = c
	defer client.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Error("read:", err)
			break
		}
		log.Info("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Error("write:", err)
			break
		}
	}
}

func Start(addr string, fn func(bool, error)) {
	Status = true
	fn(Status, nil)
	m := http.NewServeMux()
	Server = http.Server{Addr: addr, Handler: m}
	m.HandleFunc("/", connect)
	if err := Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		Status = false
		fn(Status, err)
	}
}

func Stop(fn func(bool, error)) {
	if err := Server.Shutdown(context.Background()); err != nil {
		fn(Status, err)
	} else {
		client.Close()
		Status = false
		fn(Status, nil)
	}
}
