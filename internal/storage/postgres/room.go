package postgres

import (
	"chat-server/internal/models"
	"log"
)

//creating room 
const RoomCreateQuery = 
		`INSERT INTO rooms (name)
		 VALUES ($1) RETURNING id , name`

//fetching all rooms
const GetAllRoomsQuery = `SELECT *
		 FROM roomMembers
		 WHERE username = $1
		 ORDER BY name ASC`

//fetching room by name(unique)
const GetRoomByNameQuery = `SELECT *
		 FROM rooms
		 WHERE name = $1`


func (s *Store) CreateRoom(room string) (*models.StoredRoom) {
	var result models.StoredRoom
	err:= s.db.QueryRow(RoomCreateQuery,room).Scan(&result.ID,&result.Name)
	if(err!=nil){
		log.Println(err);
		return nil;
	}
	return &result;

}

func (s *Store) GetRoomByName(room string) (*models.StoredRoom,error) {
	DBroom, err := s.db.Query(GetRoomByNameQuery,room,)
	if(err!=nil){
		log.Println(err)
		return nil,err
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
	return &r,err
}
func (s *Store) GetAllRooms() ([]*models.StoredRoom, error) {
	rows, err := s.db.Query(GetAllRoomsQuery,)
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