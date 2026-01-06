package postgres

import (
	"chat-server/internal/models"
	"log"
)

// adding user with given name(which is unique) to the room
const AddUserToRoomQuery = `INSERT INTO room_members (room_id, username, room_name) VALUES ($1, $2, $3)`

// removing user with given name(which is unique) from the room
const RemoveUserFromRoomQuery = `DELETE FROM room_members WHERE room_id = $1 AND username = $2`

// get all the users in the room with given id
const GetUsersInRoomQuery = `SELECT * FROM room_members WHERE room_id = $1`

// fetching all rooms of a user
const GetAllRoomsOfUserQuery = `SELECT room_id,room_name
		 FROM room_members
		 WHERE username=$1
		 ORDER BY username ASC`

// add user to room if err==nil means user is added
func (s *Store) AddUserToRoom(roomId int, username string, roomName string) error {
	_, err := s.db.Exec(
		AddUserToRoomQuery,
		roomId,
		username,
		roomName,
	)
	return err
}

// delete user to room if err==nil means user is added
func (s *Store) RemoveUserFromRoom(roomId int, username string) error {
	_, err := s.db.Exec(
		RemoveUserFromRoomQuery,
		roomId,
		username,
	)
	return err
}

func (s *Store) GetUsersInRoom(roomId int) ([]*models.RoomMember, error) {
	rows, err := s.db.Query(
		GetUsersInRoomQuery,
		roomId,
	)
	if err != nil {
		log.Println("Can't fetch users from room with id ", roomId)
		return nil, err
	}
	defer rows.Close()
	var members []*models.RoomMember
	for rows.Next() {
		var m models.RoomMember
		if err := rows.Scan(
			&m.RoomID,
			&m.Username,
			&m.JoinedAt,
		); err != nil {
			return nil, err
		}
		members = append(members, &m)
	}
	return members, nil
}

func (s *Store) GetRoomsOfUser(username string) ([]*models.StoredRoom, error) {
	rows, err := s.db.Query(GetAllRoomsOfUserQuery, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rooms []*models.StoredRoom
	for rows.Next() {
		var r models.StoredRoom
		if err := rows.Scan(
			&r.ID,
			&r.Name,
		); err != nil {
			return nil, err
		}
		rooms = append(rooms, &r)
	}
	return rooms, nil
}
