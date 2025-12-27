package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan Message
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

		h.Broadcast <- msg //send msg for broadcasting
		// log.Println("Message sent for broadcasting:-", string(msg))
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
