package storage

import (
	"chat-server/internal/models"
	"context"
)

type StorageInterface interface {

	//message operations

	//limit=20=>last 20 messages
	SaveMessage(ctx context.Context, msg string, roomId int, userName string) error
	GetRecentMessages(ctx context.Context, roomId int, limit int, lastid int) ([]models.Message, error)

	//dm
	SendDirectMessage(ctx context.Context, msg string, reciever string, user string) error
	GetRecentDirectMessages(ctx context.Context, username string, limit int, lastid int) ([]models.Message, error)

	//user operations
	CreateUserIfNotExists(ctx context.Context, user string)
	GetUserByName(ctx context.Context, user string) (int, error)

	//room operations
	CreateRoom(ctx context.Context, room string, name string) (*models.StoredRoom, error)
	GetRoomByName(ctx context.Context, name string) (*models.StoredRoom, error)
	GetRoomById(ctx context.Context, id int) (*models.StoredRoom, error)
	GetAllRooms(ctx context.Context) ([]*models.StoredRoom, error)

	//room-user operations
	AddUserToRoom(ctx context.Context, roomId int, username string) error
	RemoveUserFromRoom(ctx context.Context, roomId int, username string) error
	GetUsersInRoom(ctx context.Context, roomId int) ([]*models.RoomMember, error)
	GetRoomsOfUser(ctx context.Context, username string) ([]*models.StoredRoom, error)
}
