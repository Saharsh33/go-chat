package websocket

import "log"

type Hub struct {

	//string is username
	Clients map[string]*Client

	// Channels to handle server operations concurrently
	Rooms map[Room]map[*Client]bool

	//go-channels for different operations
	JoinRoom    chan RoomOps
	LeaveRoom   chan RoomOps
	CreateRoom  chan RoomOps
	Register    chan *Client
	Unregister  chan *Client
	SendMessage chan Message
}

func NewHub() *Hub {
	return &Hub{
		Clients:     make(map[string]*Client),
		Rooms:       make(map[Room]map[*Client]bool),
		JoinRoom:    make(chan RoomOps),
		LeaveRoom:   make(chan RoomOps),
		CreateRoom:  make(chan RoomOps),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		SendMessage: make(chan Message),
	}
}

func (h *Hub) Run() {

	for {

		select {

		case client := <-h.Register:

			// Registering the client by mapping in h.Clients
			h.Clients[client.Username] = client

		case client := <-h.Unregister:

			// Unregistering the client
			if _, ok := h.Clients[client.Username]; ok {
				//client exists
				delete(h.Clients, client.Username)
				for _, clients := range h.Rooms {
					delete(clients, client)
				}
				// close send channel
				close(client.Send)
			}

		case message := <-h.SendMessage:

			switch message.Type {

			case "broadcast":

				// Send message to all clients
				message.Room = "/Broadcast"
				for _, client := range h.Clients {

					select {

					case client.Send <- message:

					default:
						//if client is slow or dead
						close(client.Send)

						for _, clients := range h.Rooms {
							delete(clients, client)
						}

						delete(h.Clients, client.Username)
					}
				}

			case "messageRoom":

				// Send message to particular room
				var tempRoom Room = Room{name: message.Room}
				_, ok1 := h.Rooms[tempRoom]

				if ok1 {

					// If room exists
					_, ok2 := h.Rooms[tempRoom][h.Clients[message.User]]

					if ok2 {

						// If client exists in that room
						for client := range h.Rooms[Room{name: message.Room}] {

							//for all clients in the room
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

			// Join a room
			var tempRoom Room = Room{name: JoinRoomDetails.roomDetails}
			_, ok := h.Rooms[tempRoom]

			if ok {

				// Room exists
				h.Rooms[tempRoom][JoinRoomDetails.clientDetails] = true
				log.Printf("Joined %s succesfully\n", JoinRoomDetails.roomDetails)
			} else {
				log.Printf("%s room doesnt exis\n", JoinRoomDetails.roomDetails)
			}

		case LeaveRoomDetails := <-h.LeaveRoom:

			// Leave a room
			var tempRoom Room = Room{name: LeaveRoomDetails.roomDetails}
			_, ok := h.Rooms[tempRoom]

			if ok {

				delete(h.Rooms[tempRoom], LeaveRoomDetails.clientDetails)
				log.Printf("Left %s succesfully\n", LeaveRoomDetails.roomDetails)
			} else {
				log.Printf("Couldn't leave %s\n", LeaveRoomDetails.roomDetails)
			}

		case CreateRoomDetails := <-h.CreateRoom:

			// Create a room
			log.Println("Recieved for room creation!!")

			var tempRoom Room = Room{name: CreateRoomDetails.roomDetails}
			_, ok := h.Rooms[tempRoom]
			if ok {

				log.Println("Room already exists")
			} else {
				h.Rooms[tempRoom] = map[*Client]bool{}
				h.Rooms[tempRoom][CreateRoomDetails.clientDetails] = true

				log.Println("Room created and joined !! ", tempRoom)
			}
		}

	}
}
