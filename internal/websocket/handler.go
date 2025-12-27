package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true //returning true all the time(for now)
	},
}

func ServeWS(h *Hub, w http.ResponseWriter, r *http.Request) {

	//upgrading from http to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	//create client
	client := &Client{
		Conn: conn,
		Send: make(chan []byte),
	}

	//register client
	h.Register <- client

	//goroutine
	go func() {
		defer func() {
			h.Unregister <- client
			conn.Close()
		}()

		//temporary
		for {
			msgType, msg, err := conn.ReadMessage()
			//break the loop upon error
			if err != nil {
				log.Println(err)
				break
			}
			log.Println(msgType, msg)
		}
	}()

}
