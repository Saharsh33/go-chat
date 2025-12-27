package websocket

import "log"

type Hub struct {
	Clients    map[*Client]bool
	Rooms      map[Room]map[*Client]bool
	JoinRoom   chan *RoomOps
	LeaveRoom  chan *RoomOps
	CreateRoom chan *RoomOps
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:

			h.Clients[client] = true

		case client := <-h.Unregister:

			_, exists := h.Clients[client] //break to key,value pair
			if exists {                    //if client exists then remove from map and close channel
				delete(h.Clients, client)
				close(client.Send)
			}

		case message := <-h.Broadcast: //broadcast msg is sent
			switch message.Type {
			case "broadcast":
				for client := range h.Clients { //for all clients
					select {
					case client.Send <- message:
					default:
						close(client.Send) //if client is slow or dead
						delete(h.Clients, client)
					}
				}
			case "messageRoom":
				for client := range h.Rooms[Room{name: message.Room}] { //for all clients
					select {
					case client.Send <- message:
					default:
						close(client.Send) //if client is slow or dead
						delete(h.Clients, client)
					}
				}

			}
		case JoinRoomDetails := <-h.JoinRoom:
			var tempRoom Room = Room{name: JoinRoomDetails.roomDetails}
			_, ok := h.Rooms[tempRoom]
			if ok {
				h.Rooms[tempRoom][JoinRoomDetails.clientDetails] = true
			}
		case LeaveRoomDetails := <-h.LeaveRoom:
			var tempRoom Room = Room{name: LeaveRoomDetails.roomDetails}
			_, ok := h.Rooms[tempRoom]
			if ok {
				h.Rooms[tempRoom][LeaveRoomDetails.clientDetails] = false
			}

		case CreateRoomDetails := <-h.CreateRoom:
			var tempRoom Room = Room{name: CreateRoomDetails.roomDetails}
			_, ok := h.Rooms[tempRoom]
			if ok {
				log.Println("Room already exists")
			} else {
				h.Rooms[tempRoom] = map[*Client]bool{}
			}
		}

	}
}
