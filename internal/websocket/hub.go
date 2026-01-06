package websocket

import (
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
	store       storage.StorageInterface
}

func NewHub(store storage.StorageInterface) *Hub {
	return &Hub{
		Clients:     make(map[string]*Client),
		Rooms:       make(map[Room]map[*Client]bool),
		JoinRoom:    make(chan RoomOps),
		LeaveRoom:   make(chan RoomOps),
		CreateRoom:  make(chan RoomOps),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		SendMessage: make(chan Message),
		store:       store,
	}
}

func (h *Hub) Run() {
	log.Println("Preparing HUB!!")

	log.Println("HUB running!")
	for {

		select {

		case client := <-h.Register:

			// Registering the client by mapping in h.Clients
			h.Clients[client.Username] = client
			roomsOfUser, err := h.store.GetRoomsOfUser(client.Username)
			if err != nil {
				log.Println("Can't fetch user's joined room details ", err)
			} else {
				for _, room := range roomsOfUser {
					if (h.Rooms[Room{name: room.Name}] == nil) {
						h.Rooms[Room{name: room.Name}] = map[*Client]bool{}
					}
					h.Rooms[Room{name: room.Name}][client] = true
				}
				client.Send <- Message{Type: MsgSystem, User: "system", Room: "system", Content: client.Username + " Registered Successfully"}
			}

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

			room, ok2 := h.store.GetRoomByName(tempRoom.name)

			if ok2 != nil {

				log.Println("Unable to find room!!(DB Query error)", ok2)

			} else {

				if room != nil {

					err := h.store.AddUserToRoom(room.ID, JoinRoomDetails.clientDetails.Username)

					if err == nil {

						if h.Rooms[tempRoom] == nil {
							h.Rooms[tempRoom] = map[*Client]bool{}
						}

						h.Rooms[tempRoom][JoinRoomDetails.clientDetails] = true
						JoinRoomDetails.clientDetails.Send <- Message{Type: MsgSystem, User: "system", Room: "system", Content: "Joined " + JoinRoomDetails.roomDetails + " Succesfully!!"}
						log.Println(JoinRoomDetails.clientDetails.Username, " joined ", JoinRoomDetails.roomDetails, " succesfully")

					} else {

						log.Println("DB Query error : Couldn't add ", JoinRoomDetails.clientDetails.Username, " to ", JoinRoomDetails.roomDetails)

					}

				} else {

					log.Println("No such room found!!")

				}
			}

		case LeaveRoomDetails := <-h.LeaveRoom:

			// Leave a room
			var tempRoom Room = Room{name: LeaveRoomDetails.roomDetails}
			_, ok := h.Rooms[tempRoom]

			if ok {

				room, ok2 := h.store.GetRoomByName(tempRoom.name)

				if ok2 != nil {

					log.Println("DB error room not found 102")

				} else {

					err := h.store.RemoveUserFromRoom(room.ID, LeaveRoomDetails.clientDetails.Username)

					if err != nil {

						log.Println("DB error code 101")

					} else {

						delete(h.Rooms[tempRoom], LeaveRoomDetails.clientDetails)

						LeaveRoomDetails.clientDetails.Send <- Message{Type: MsgSystem, User: "system", Room: "system", Content: "Left " + LeaveRoomDetails.roomDetails + " Succesfully!!"}
						log.Printf("Left %s succesfully\n", LeaveRoomDetails.roomDetails)
					}
				}
			} else {

				LeaveRoomDetails.clientDetails.Send <- Message{Type: MsgSystem, User: "system", Room: "system", Content: "Unable to leave " + LeaveRoomDetails.roomDetails}
				log.Printf("Couldn't leave %s\n", LeaveRoomDetails.roomDetails)
			}

		case CreateRoomDetails := <-h.CreateRoom:

			// Create a room
			log.Println("Recieved for room creation!!")

			var tempRoom Room = Room{name: CreateRoomDetails.roomDetails}
			_, ok1 := h.Rooms[tempRoom]

			if ok1 {

				CreateRoomDetails.clientDetails.Send <- Message{Type: MsgSystem, User: "system", Room: "system", Content: CreateRoomDetails.roomDetails + " already exists!!"}
				log.Println("Room already exists")

			} else {

				if _, ok2 := h.store.GetRoomByName(tempRoom.name); ok2 != nil {

					CreateRoomDetails.clientDetails.Send <- Message{Type: MsgSystem, User: "system", Room: "system", Content: CreateRoomDetails.roomDetails + " already exists!!"}
					log.Println("Room already exists")

				} else {

					result, err := h.store.CreateRoom(CreateRoomDetails.roomDetails, CreateRoomDetails.clientDetails.Username)
					if err != nil {

						log.Println(err)

					} else {

						tempRoom = Room{name: result.Name}
						h.Rooms[tempRoom] = map[*Client]bool{}
						h.Rooms[tempRoom][CreateRoomDetails.clientDetails] = true

						CreateRoomDetails.clientDetails.Send <- Message{Type: MsgSystem, User: "system", Room: "system", Content: "Created  " + CreateRoomDetails.roomDetails + " Succesfully!!"}
						log.Println("Room created and joined !! ", tempRoom)

					}
				}
			}
		}
	}
}
