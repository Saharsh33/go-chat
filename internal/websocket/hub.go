package websocket

import "log"

type Hub struct {
	Clients    map[string]*Client //string is username
	Rooms      map[Room]map[*Client]bool
	JoinRoom   chan RoomOps
	LeaveRoom  chan RoomOps
	CreateRoom chan RoomOps
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Rooms:      make(map[Room]map[*Client]bool),
		JoinRoom:   make(chan RoomOps),
		LeaveRoom:  make(chan RoomOps),
		CreateRoom: make(chan RoomOps),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:

			h.Clients[client.Username] = client

		case client := <-h.Unregister:
			if _, ok := h.Clients[client.Username]; ok {
				delete(h.Clients, client.Username)

				for _, clients := range h.Rooms {
					delete(clients, client)
				}

				close(client.Send)
			}

		case message := <-h.Broadcast: //broadcast msg is sent
			switch message.Type {
			case "broadcast":
				message.Room = "From Broadcast Channel"
				for _, client := range h.Clients { //for all clients
					select {
					case client.Send <- message:
					default:
						close(client.Send) //if client is slow or dead
						for _, clients := range h.Rooms {
							delete(clients, client)
						}
						delete(h.Clients, client.Username)
					}
				}
			case "messageRoom":
				var tempRoom Room = Room{name: message.Room}
				_, ok1 := h.Rooms[tempRoom]

				if ok1 {
					_, ok2 := h.Rooms[tempRoom][h.Clients[message.User]]
					if ok2 {
						for client := range h.Rooms[Room{name: message.Room}] { //for all clients in the room
							select {
							case client.Send <- message:
							default:
								close(client.Send) //if client is slow or dead
								delete(h.Clients, client.Username)
							}
						}
					} else {
						log.Println(message.User, " is not in the room:-", tempRoom.name)
					}
				} else {
					log.Println("Room doesn't exists or you are not in the room!! ", tempRoom)
				}

			}
		case JoinRoomDetails := <-h.JoinRoom:
			var tempRoom Room = Room{name: JoinRoomDetails.roomDetails}
			_, ok := h.Rooms[tempRoom]
			if ok {
				h.Rooms[tempRoom][JoinRoomDetails.clientDetails] = true
				log.Printf("Joined %s succesfully\n", JoinRoomDetails.roomDetails)
			} else {
				log.Printf("Couldn't join %s\n", JoinRoomDetails.roomDetails)
			}
		case LeaveRoomDetails := <-h.LeaveRoom:
			var tempRoom Room = Room{name: LeaveRoomDetails.roomDetails}
			_, ok := h.Rooms[tempRoom]
			if ok {
				delete(h.Rooms[tempRoom], LeaveRoomDetails.clientDetails)
				log.Printf("Left %s succesfully\n", LeaveRoomDetails.roomDetails)
			} else {
				log.Printf("Couldn't leave %s\n", LeaveRoomDetails.roomDetails)

			}

		case CreateRoomDetails := <-h.CreateRoom:
			log.Println("Recieved for room creation!!")
			var tempRoom Room = Room{name: CreateRoomDetails.roomDetails}
			_, ok := h.Rooms[tempRoom]
			if ok {
				log.Println("Room already exists")
			} else {
				h.Rooms[tempRoom] = map[*Client]bool{}
				h.Rooms[tempRoom][CreateRoomDetails.clientDetails] = true
				log.Println("Room created!! ", tempRoom)
			}
		}

	}
}
