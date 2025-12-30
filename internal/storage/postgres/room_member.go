package postgres

import (
	"chat-server/internal/models"
	"log"
)

//add user to room if err==nil means user is added
func (s *Store) AddUserToRoom(roomId int , username string) error{
	_, err := s.db.Exec(
		`INSERT INTO roomMembers (room_id, username)
		 VALUES ($1, $2)`,
		roomId,
		username,
	)
	return err
}

//delete user to room if err==nil means user is added
func (s *Store) RemoveUserFromRoom(roomId int, username string) error{
	_, err := s.db.Exec(
		`DELETE FROM roomMembers 
		WHERE room_id=$1 AND username=$2`,
		roomId,
		username,
	)
	return err
}


func (s *Store) GetUsersInRoom(roomId int) ([]*models.RoomMember, error){
	rows, err := s.db.Query(
		`SELECT *
		 FROM roomMembers
		 WHERE room_id=$1`,
		 roomId,
	)
	if err != nil {
		log.Println("Can't fetch users from room with id ",roomId)
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
