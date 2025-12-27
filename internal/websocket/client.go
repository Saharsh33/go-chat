package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Send     chan Message
	Username string
}

// readpump can detect disconnection error
func (c *Client) readPump(h *Hub) {
	defer func() {
		h.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, raw, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		var msg Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			log.Println("invalid json:", err)
			continue
		}
		msg.User = c.Username
		var roomOpsDetails RoomOps = RoomOps{clientDetails: c, roomDetails: msg.Room}

		switch msg.Type {
		case "join":
			h.JoinRoom <- roomOpsDetails
		case "leave":
			h.LeaveRoom <- roomOpsDetails
		case "create":
			log.Printf("Sent for room creation!! %+v", roomOpsDetails)
			h.CreateRoom <- roomOpsDetails
		case "messageRoom":
			h.Broadcast <- msg
		case "broadcast":
			h.Broadcast <- msg //send msg for broadcasting
			// log.Println("Message sent for broadcasting:-", string(msg))
		default:
			fmt.Println("invalid Request:")
		}

	}
}

func (c *Client) writePump() {
	defer c.Conn.Close()

	for message := range c.Send {
		data, err := json.Marshal(message)
		if err != nil {
			log.Println("marshal error:", err)
			continue
		}
		err = c.Conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println(err)
			break
		}
		// log.Println("Message writen:-", string(message))
	}
}
