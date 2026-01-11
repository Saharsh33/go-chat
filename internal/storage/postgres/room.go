package postgres

import (
	"chat-server/internal/models"
	"log"
)

// creating room
const RoomCreateQuery = `INSERT INTO rooms (name, description, created_by)
		 VALUES ($1, $2, $3) RETURNING id , name`

// fetching all rooms
const GetAllRoomsQuery = `SELECT *
		 FROM roomMembers
		 WHERE username = $1
		 ORDER BY name ASC`

// fetching room by name(unique)
const GetRoomByNameQuery = `SELECT id,name
		 FROM rooms
		 WHERE name = $1`

const GetRoomByIdQuery = `SELECT id,name
		 FROM rooms
		 WHERE id = $1`

func (s *Store) CreateRoom(room string, name string) (*models.StoredRoom, error) {
	var result models.StoredRoom
	roomdescription := ""
	err := s.db.QueryRow(RoomCreateQuery, room, roomdescription, name).Scan(&result.ID, &result.Name)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = s.AddUserToRoom(result.ID, name)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &result, err

}

func (s *Store) GetRoomByName(room string) (*models.StoredRoom, error) {
	DBroom, err := s.db.Query(GetRoomByNameQuery, room)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var r models.StoredRoom
	for DBroom.Next() {
		if err := DBroom.Scan(
			&r.ID,
			&r.Name,
		); err != nil {
			return nil, err
		}
	}
	return &r, err
}
func (s *Store) GetRoomById(id int) (*models.StoredRoom, error) {
	DBroom, err := s.db.Query(GetRoomByIdQuery, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var r models.StoredRoom
	for DBroom.Next() {
		if err := DBroom.Scan(
			&r.ID,
			&r.Name,
		); err != nil {
			return nil, err
		}
	}
	return &r, err
}
func (s *Store) GetAllRooms() ([]*models.StoredRoom, error) {
	rows, err := s.db.Query(GetAllRoomsQuery)
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
