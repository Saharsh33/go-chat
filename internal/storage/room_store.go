package storage

import "chat-server/internal/models"

type RoomStore interface {
	CreateRoom(name string) error
	GetAllRooms() ([]models.StoredRoom, error)
	GetRoomByName(name string) (bool, error)

	AddUserToRoom(roomId int, username string) error
	RemoveUserFromRoom(roomName, username string) error
	GetUsersInRoom(roomName string) ([]string, error)
}
