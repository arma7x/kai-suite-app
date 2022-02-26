package websocketserver

import (
	"flag"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	addr = flag.String("addr", "192.168.43.33:4040", "http service address")
)

var client *websocket.Conn

func connect(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	client = c
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func Start() {
	log.Print("Vroom, Vroom")
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", connect)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
