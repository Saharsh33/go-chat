package main

import (
	"chat-server/internal/websocket"
	"net/http"
)

func main() {
	hub := websocket.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWS(hub, w, r)
	})

	http.ListenAndServe(":3000", nil)
}
