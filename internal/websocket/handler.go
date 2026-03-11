package websocket

import (
	"context"
	"log"
	"net/http"

	"chat-server/internal/auth"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true //returning true all the time(for now)
	},
}

func ServeWS(h *Hub, w http.ResponseWriter, r *http.Request, jwtSecret string) {
	// Extract and validate the JWT token from the query parameter.
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	username, err := auth.ValidateToken(tokenString, jwtSecret)
	if err != nil {
		log.Println("invalid token:", err)
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	// Upgrade HTTP connection to WebSocket.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	ctx, cancel := context.WithCancel(r.Context())

	//create client
	client := &Client{
		Conn:     conn,
		Username: username,
		Send:     make(chan Message),
		Ctx:      ctx,
		Cancel:   cancel,
	}

	//register client
	h.Register <- client

	// Reading and writing msgs to client till connection closes
	go client.writePump()
	go client.readPump(h)
}

