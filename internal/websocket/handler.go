package websocket

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true //returning true all the time(for now)
	},
}

var temp int = 0

func ServeWS(h *Hub, w http.ResponseWriter, r *http.Request) {

	// This function gets called when client send http req to server
	// upgrading from http to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	//create client
	client := &Client{
		Conn:     conn,
		Username: "user-" + strconv.Itoa(temp%3),
		Send:     make(chan Message),
	}

	temp++

	//register client
	h.Register <- client

	// Reading and writing msgs to client till connection closes
	go client.writePump()
	go client.readPump(h)
}
