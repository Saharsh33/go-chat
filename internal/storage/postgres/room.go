package postgres

import (
	"chat-server/internal/models"
	"log"
)

func (s *Store) CreateRoom(room *models.StoredRoom) error {
	_, err := s.db.Exec(
		`INSERT INTO rooms (id, name, created_at)
		 VALUES ($1, $2, $3)`,
		room.ID,
		room.Name,
		room.CreatedAt,
	)
	return err
}
func (s *Store) GetAllRooms(room models.StoredRoom) ([]*models.StoredRoom,error) {
	rows, err := s.db.Query(
		`SELECT id, name, created_at
		 FROM rooms
		 ORDER BY createdAt DESC`,
	)
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
			&r.CreatedAt,
		); err != nil {
			return nil, err
		}
		rooms = append(rooms, &r)
	}
	return rooms, nil
}
func (s *Store) GetRoomByName(room string) (*models.StoredRoom,error) {
	DBroom, err := s.db.Query(
		`SELECT *
		 FROM rooms
		 WHERE name = $1`,
		room,
	)
	if(err!=nil){
		log.Println(err)
		return nil,err
	}
	var r models.StoredRoom
	for DBroom.Next() {
		if err := DBroom.Scan(
			&r.ID,
			&r.Name,
			&r.CreatedAt,
		); err != nil {
			return nil, err
		}
	}
	return &r,err
}
