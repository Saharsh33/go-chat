package websocket

import (
	"chat-server/internal/storage"
	"context"
	"log"
	"strconv"
	"time"
)

const LIMIT int = 20

type Hub struct {

	//username to client struct
	Clients map[string]*Client

	Rooms map[int]Room

	// room to client
	UsersOfRoom map[int]map[string]struct{} // int to username   try to make it int to int
	// client to room
	RoomsOfUser map[string]map[int]struct{}
	// Channels to handle server operations concurrently

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
		UsersOfRoom: make(map[int]map[string]struct{}),
		RoomsOfUser: make(map[string]map[int]struct{}),
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
	log.Println("HUB running!")

	for {

		select {

		case client := <-h.Register:
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			h.store.CreateUserIfNotExists(ctx, client.Username) // Should be handled somewhere else

			// Registering the client by mapping in h.Clients
			h.Clients[client.Username] = client

			//loading dms

			messages, err := h.store.GetRecentDirectMessages(ctx, client.Username, LIMIT, 0)

			if err != nil {
				log.Println("code : 105", err)
			} else {
				for _, directMessageOfUser := range messages {
					client.Send <- Message{Type: MsgDirectMessage, User: directMessageOfUser.User, Receiver: directMessageOfUser.Receiver, Content: directMessageOfUser.Content}
				}
			}

			h.RoomsOfUser[client.Username] = map[int]struct{}{}
			DBroomsOfUser, err := h.store.GetRoomsOfUser(ctx, client.Username)
			if err != nil {
				log.Println("Can't fetch user's joined room details ", err)
			} else {
				for _, room := range DBroomsOfUser {
					if h.UsersOfRoom[room.ID] == nil {
						h.UsersOfRoom[room.ID] = map[string]struct{}{}
						log.Println(room.ID, "included in map")
					}
					h.UsersOfRoom[room.ID][client.Username] = struct{}{} //change the structure of the map
					messages, err := h.store.GetRecentMessages(ctx, room.ID, LIMIT, 0)
					if err != nil {
						log.Println("code : 103", err)
					} else {
						for _, messagesOfRoom := range messages {
							client.Send <- Message{Type: MsgRoomMessage, User: messagesOfRoom.User, Room: messagesOfRoom.Room, Content: messagesOfRoom.Content}
						}
					}
					client.Send <- Message{Type: MsgSystem, User: "system", Room: -1, Content: client.Username + " Registered Successfully"}
				}
				client.Send <- Message{Type: MsgSystem, User: "system", Room: -1, Content: client.Username + " Registered Successfully"}
			}

		case client := <-h.Unregister:

			// Unregistering the client
			if _, ok := h.Clients[client.Username]; ok {
				//client exists
				delete(h.Clients, client.Username)
				for room, _ := range h.RoomsOfUser[client.Username] {
					delete(h.UsersOfRoom[room], client.Username)
					// check if room is empty
				}
				delete(h.RoomsOfUser, client.Username)
				client.Send <- Message{Type: MsgSystem, User: "system", Room: -1, Content: client.Username + " Unregistered Successfully"}
				// close send channel
				close(client.Send)
			} else {
				log.Println("Client doesn't exist in map")
			}
		case message := <-h.SendMessage:

			switch message.Type {

			case MsgBroadcast:

				// Send message to all clients
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				err := h.store.SaveMessage(ctx, message.Content, message.Room, message.User)
				if err != nil {
					log.Println("Can't Save message!! ", err)
					break
				}

				message.Room = 0
				for _, client := range h.Clients {
					select {

					case client.Send <- message:

					default:
						//if client is slow or dead
						close(client.Send)

						delete(h.Clients, client.Username)
						for room, _ := range h.RoomsOfUser[client.Username] {
							delete(h.UsersOfRoom[room], client.Username)
							// check if room is empty
						}
						delete(h.RoomsOfUser, client.Username)
					}
				}
				//h.Store.SaveMessage(models.Message{ID:,Type:,Room:,User:,Content:,CreatedAt:})
				h.Clients[message.User].Send <- Message{Type: MsgSystem, User: "system", Room: -1, Content: "Message sent to broadcast"}

			case MsgRoomMessage:
				roomId := message.Room
				// Send message to particular room

				_, ok1 := h.UsersOfRoom[roomId]

				if ok1 {

					// If room exists
					_, ok2 := h.UsersOfRoom[roomId][message.User]

					if ok2 {

						// If client exists in that room

						ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
						defer cancel()
						err := h.store.SaveMessage(ctx, message.Content, message.Room, message.User)
						if err != nil {
							log.Println("Can't Save message!! ", err)
							break
						}
						for user := range h.UsersOfRoom[roomId] {

							//for all clients in the room
							client := h.Clients[user]
							select {

							case client.Send <- message:

							default:
								close(client.Send)

								delete(h.Clients, client.Username)
								for room, _ := range h.RoomsOfUser[client.Username] {
									delete(h.UsersOfRoom[room], client.Username)
									// check if room is empty
								}
								delete(h.RoomsOfUser, client.Username)
							}
						}
						h.Clients[message.User].Send <- Message{Type: MsgSystem, User: "system", Room: -1, Content: "Message sent"}

					} else {
						h.Clients[message.User].Send <- Message{Type: MsgSystem, User: "system", Room: -1, Content: "User not part of room"}
						log.Println(message.User, " is not in the room:-", message.Room)
					}
				} else {
					h.Clients[message.User].Send <- Message{Type: MsgSystem, User: "system", Room: -1, Content: "Room " + strconv.Itoa(message.Room) + " doesn't exist"}
					log.Println("Room doesn't exists", message.Room)
				}

			case MsgDirectMessage:
				// send direct msg
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()

				_, err := h.store.GetUserByName(ctx, message.Receiver)
				if err != nil {
					log.Println("User not found", err)
					break
				}
				err = h.store.SendDirectMessage(ctx, message.Content, message.Receiver, message.User)

				if err != nil {
					log.Println("hub.go/Run/SendMessage/DirectMessage", err)
				} else {
					h.Clients[message.User].Send <- Message{Type: MsgSystem, User: message.User, Receiver: message.Receiver, Content: message.Content}
					h.Clients[message.Receiver].Send <- Message{Type: MsgSystem, User: message.User, Receiver: message.Receiver, Content: message.Content}
				}
			case MsgNextRoomMessages:

				lastid, _ := strconv.Atoi(message.Content)
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				msgs, err := h.store.GetRecentMessages(ctx, message.Room, LIMIT, lastid)
				if err != nil {
					log.Println("Error retrieving next msgs", err)
				} else {
					for _, messagesOfRoom := range msgs {
						h.Clients[message.User].Send <- Message{Type: MsgRoomMessage, User: messagesOfRoom.User, Room: messagesOfRoom.Room, Content: messagesOfRoom.Content}
					}

				}

			case MsgNextDirectMessages:
				lastid, _ := strconv.Atoi(message.Content)
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				msgs, err := h.store.GetRecentDirectMessages(ctx, message.User, LIMIT, lastid)
				if err != nil {
					log.Println("Error retrieving next msgs", err)
				} else {
					for _, directMessageOfUser := range msgs {
						h.Clients[message.User].Send <- Message{Type: MsgDirectMessage, User: directMessageOfUser.User, Receiver: directMessageOfUser.Receiver, Content: directMessageOfUser.Content}
					}

				}

			}

		case JoinRoomDetails := <-h.JoinRoom:

			client := JoinRoomDetails.clientDetails
			roomID := JoinRoomDetails.roomDetails

			// Per-operation timeout context
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			// Join a room
			flag := false
			_, ok2 := h.UsersOfRoom[roomID]

			if ok2 != true {

				log.Println("Unable to find room in map")

				_, ok1 := h.store.GetRoomById(ctx, roomID)
				if ok1 == nil {
					flag = true
					h.UsersOfRoom[roomID] = map[string]struct{}{}

				} else {
					log.Println("No such room found in db!!", ok1)
				}

			} else {
				flag = true
			}
			if flag == true {

				err := h.store.AddUserToRoom(ctx, roomID, client.Username)

				if err == nil {

					h.UsersOfRoom[roomID][client.Username] = struct{}{}
					h.RoomsOfUser[client.Username][roomID] = struct{}{}
					JoinRoomDetails.clientDetails.Send <- Message{Type: MsgSystem, User: "system", Room: -1, Content: "Joined " + strconv.Itoa(JoinRoomDetails.roomDetails) + " Succesfully!!"}
					log.Println(JoinRoomDetails.clientDetails.Username, " joined ", JoinRoomDetails.roomDetails, " succesfully")

				} else {

					log.Println("DB Query error : Couldn't add ", JoinRoomDetails.clientDetails.Username, " to ", JoinRoomDetails.roomDetails, err)

				}

			} else {

				log.Println("No such room found!!")

			}

		case LeaveRoomDetails := <-h.LeaveRoom:

			client := LeaveRoomDetails.clientDetails
			roomID := LeaveRoomDetails.roomDetails
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			// Leave a room
			_, ok := h.UsersOfRoom[roomID]

			if ok {

				err := h.store.RemoveUserFromRoom(ctx, roomID, client.Username)

				if err != nil {

					log.Println("DB error code 101")

				} else {

					delete(h.UsersOfRoom[roomID], client.Username)
					delete(h.RoomsOfUser[client.Username], roomID)

					LeaveRoomDetails.clientDetails.Send <- Message{Type: MsgSystem, User: "system", Room: -1, Content: "Left " + strconv.Itoa(roomID) + " Succesfully!!"}
					log.Printf("Left %d succesfully\n", roomID)
				}

			} else {

				LeaveRoomDetails.clientDetails.Send <- Message{Type: MsgSystem, User: "system", Room: -1, Content: "Unable to leave " + strconv.Itoa(LeaveRoomDetails.roomDetails)}
				log.Printf("Couldn't leave %d\n", LeaveRoomDetails.roomDetails)
			}

		case CreateRoomDetails := <-h.CreateRoom:

			// Create a room
			client := CreateRoomDetails.clientDetails
			roomName := CreateRoomDetails.roomName
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			log.Println("Recieved for room creation!!")

			if _, ok2 := h.store.GetRoomByName(ctx, roomName); ok2 != nil {

				client.Send <- Message{Type: MsgSystem, User: "system", Room: -1, Content: strconv.Itoa(CreateRoomDetails.roomDetails) + " already exists!!"}
				log.Println("Room already exists")

			} else {

				result, err := h.store.CreateRoom(ctx, roomName, client.Username)

				if err != nil {

					log.Println(err)

				} else {

					h.UsersOfRoom[result.ID] = map[string]struct{}{}
					h.UsersOfRoom[result.ID][client.Username] = struct{}{}
					h.RoomsOfUser[client.Username][result.ID] = struct{}{}

					CreateRoomDetails.clientDetails.Send <- Message{Type: MsgSystem, User: "system", Room: -1, Content: "Created  " + strconv.Itoa(CreateRoomDetails.roomDetails) + " Succesfully!!"}
					log.Println("Room created and joined !! ", result.Name)

				}
			}

		}
	}
}
