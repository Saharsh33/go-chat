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

	//cleanup after disconnection
	defer func() {
		h.Unregister <- c
		c.Conn.Close()
	}()

	//readpump logic
	for {
		_, raw, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		//check if valid JSON
		var msg Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			log.Println("invalid json:", err)
			continue
		}
		msg.User = c.Username
		roomId := msg.Room

		//different cases of message types
		switch msg.Type {

		//join room
		case MsgJoinRoom:
			// room name will be null and not considered
			roomOpsDetails := RoomOps{clientDetails: c, roomDetails: roomId, roomName: msg.Content}
			h.JoinRoom <- roomOpsDetails

		//leave room
		case MsgLeaveRoom:
			// room name will be null and not considered
			roomOpsDetails := RoomOps{clientDetails: c, roomDetails: roomId, roomName: msg.Content}
			h.LeaveRoom <- roomOpsDetails

		//create room
		case MsgCreateRoom:
			// room id will be null and not considered
			roomOpsDetails := RoomOps{clientDetails: c, roomDetails: roomId, roomName: msg.Content}
			log.Printf("Sent for room creation!! %+v", roomOpsDetails)
			h.CreateRoom <- roomOpsDetails

		//message in a particular room
		case MsgRoomMessage:
			h.SendMessage <- msg

		//send to broadcast channel
		case MsgBroadcast:
			h.SendMessage <- msg
		
		case MsgDirectMessage:
			h.SendMessage <- msg

		default:
			fmt.Println("invalid Request:")
		}

	}
}

func (c *Client) writePump() {

	//cleanup after disconnection
	defer c.Conn.Close()

	//writepump logic
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
