package storage

import "chat-server/internal/models"

type StorageInterface interface {

	//message operations

	//limit=20=>last 20 messages
	SaveMessage(msg string, roomId int, userName string) error
	GetRecentMessages(roomId int, limit int) ([]models.Message, error)

	//dm
	SendDirectMessage(msg string, reciever string, user string) error

	//user operations
	CreateUserIfNotExists(user string)

	//room operations
	CreateRoom(room string, name string) (*models.StoredRoom, error)
	GetRoomByName(name string) (*models.StoredRoom, error)
	GetRoomById(id int) (*models.StoredRoom, error)
	GetAllRooms() ([]*models.StoredRoom, error)

	//room-user operations
	AddUserToRoom(roomId int, username string) error
	RemoveUserFromRoom(roomId int, username string) error
	GetUsersInRoom(roomId int) ([]*models.RoomMember, error)
	GetRoomsOfUser(username string) ([]*models.StoredRoom, error)
}
