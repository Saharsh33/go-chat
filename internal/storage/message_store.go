package storage

import "chat-server/internal/models"

type StorageInterface interface {
	SaveMessage(msg models.Message) error

	//limit=20=>last 20 messages
	GetRecentMessages(room string, limit int) ([]models.Message, error)

	//user operations
	CreateUserIfNotExists(user string)

	//room operations
	CreateRoom(room string, name string) (*models.StoredRoom, error)
	GetRoomByName(name string) (*models.StoredRoom, error)
	GetAllRooms() ([]*models.StoredRoom, error)

	//room-user operations
	AddUserToRoom(roomId int, username string, roomName string) error
	RemoveUserFromRoom(roomId int, username string) error
	GetUsersInRoom(roomId int) ([]*models.RoomMember, error)
	GetRoomsOfUser(username string) ([]*models.StoredRoom, error)
}
