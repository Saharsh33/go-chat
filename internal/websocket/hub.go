package websocket

import (
	"chat-server/internal/models"
	"chat-server/internal/storage"
	"log"
)
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
	Store 		storage.MessageStore
}

func NewHub(store storage.MessageStore) *Hub {
	return &Hub{
		Clients:     make(map[string]*Client),
		Rooms:       make(map[Room]map[*Client]bool),
		JoinRoom:    make(chan RoomOps),
		LeaveRoom:   make(chan RoomOps),
		CreateRoom:  make(chan RoomOps),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		SendMessage: make(chan Message),
		Store: 		 store,
	}
}

func (h *Hub) Run() {
	// Initial fetching
	
	for {

		select {

		case client := <-h.Register:

			// Registering the client by mapping in h.Clients
			h.Clients[client.Username] = client
			client.Send <- Message{Type: MsgSystem, User: "system", Room: "system", Content: client.Username + " Registered Successfully"}

		case client := <-h.Unregister:

			// Unregistering the client
			if _, ok := h.Clients[client.Username]; ok {
				//client exists
				delete(h.Clients, client.Username)
				for _, clients := range h.Rooms {
					delete(clients, client)
				}
				client.Send <- Message{Type: MsgSystem, User: "system", Room: "system", Content: client.Username + " Unregistered Successfully"}
				// close send channel
				close(client.Send)
			}

		case message := <-h.SendMessage:

			switch message.Type {

			case MsgBroadcast:

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
				//h.Store.SaveMessage(models.Message{ID:,Type:,Room:,User:,Content:,CreatedAt:})
				h.Clients[message.User].Send <- Message{Type: MsgSystem, User: "system", Room: "system", Content: "Message sent to broadcast"}

			case MsgRoomMessage:

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
						h.Clients[message.User].Send <- Message{Type: MsgSystem, User: "system", Room: "system", Content: "Message sent"}
					} else {
						h.Clients[message.User].Send <- Message{Type: MsgSystem, User: "system", Room: "system", Content: "User not part of room"}
						log.Println(message.User, " is not in the room:-", tempRoom.name)
					}
				} else {
					h.Clients[message.User].Send <- Message{Type: MsgSystem, User: "system", Room: "system", Content: "Room " + message.Room + " doesn't exist"}
					log.Println("Room doesn't exists", tempRoom)
				}

			}
		case JoinRoomDetails := <-h.JoinRoom:

			// Join a room
			var tempRoom Room = Room{name: JoinRoomDetails.roomDetails}
			_, ok := h.Rooms[tempRoom]

			if ok {

				// Room exists
				h.Rooms[tempRoom][JoinRoomDetails.clientDetails] = true
				JoinRoomDetails.clientDetails.Send <- Message{Type: MsgSystem, User: "system", Room: "system", Content: "Joined " + JoinRoomDetails.roomDetails + " Succesfully!!"}
				log.Println(JoinRoomDetails.clientDetails.Username, " joined ", JoinRoomDetails.roomDetails, " succesfully")
			} else {
				JoinRoomDetails.clientDetails.Send <- Message{Type: MsgSystem, User: "system", Room: "system", Content: "Room " + JoinRoomDetails.roomDetails + " doesn't exists!!"}
				log.Printf("%s room doesnt exis\n", JoinRoomDetails.roomDetails)
			}

		case LeaveRoomDetails := <-h.LeaveRoom:

			// Leave a room
			var tempRoom Room = Room{name: LeaveRoomDetails.roomDetails}
			_, ok := h.Rooms[tempRoom]

			if ok {

				delete(h.Rooms[tempRoom], LeaveRoomDetails.clientDetails)
				LeaveRoomDetails.clientDetails.Send <- Message{Type: MsgSystem, User: "system", Room: "system", Content: "Left " + LeaveRoomDetails.roomDetails + " Succesfully!!"}
				log.Printf("Left %s succesfully\n", LeaveRoomDetails.roomDetails)
			} else {
				LeaveRoomDetails.clientDetails.Send <- Message{Type: MsgSystem, User: "system", Room: "system", Content: "Unable to leave " + LeaveRoomDetails.roomDetails}
				log.Printf("Couldn't leave %s\n", LeaveRoomDetails.roomDetails)
			}

		case CreateRoomDetails := <-h.CreateRoom:

			// Create a room
			log.Println("Recieved for room creation!!")

			var tempRoom Room = Room{name: CreateRoomDetails.roomDetails}
			_, ok := h.Rooms[tempRoom]
			if ok {
				CreateRoomDetails.clientDetails.Send <- Message{Type: MsgSystem, User: "system", Room: "system", Content: CreateRoomDetails.roomDetails + " already exists!!"}
				log.Println("Room already exists")
			} else {
				h.Rooms[tempRoom] = map[*Client]bool{}
				h.Rooms[tempRoom][CreateRoomDetails.clientDetails] = true

				CreateRoomDetails.clientDetails.Send <- Message{Type: MsgSystem, User: "system", Room: "system", Content: "Created  " + CreateRoomDetails.roomDetails + " Succesfully!!"}
				log.Println("Room created and joined !! ", tempRoom)
			}
		}
	}
}
