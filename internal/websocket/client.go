package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

// readpump can detect disconnection error
func (c *Client) readPump(h *Hub) {
	defer func() {
		h.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		h.Broadcast <- msg //send msg for broadcasting
		// log.Println("Message sent for broadcasting:-", string(msg))
	}
}

func (c *Client) writePump() {
	defer c.Conn.Close()

	for message := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err)
			break
		}
		// log.Println("Message writen:-", string(message))
	}
}
