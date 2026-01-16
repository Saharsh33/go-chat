package postgres

import (
	"chat-server/internal/models"
	"database/sql"
)

const SaveMessageQuery = `INSERT INTO roommessages (room_id, username, content)
		 VALUES ($1, $2, $3)`

const GetRecentMessagesQuery = `SELECT id, room_id, username, content, created_at
		 FROM roommessages
		 WHERE room_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2`
const GetRangedMessagesQuery = `SELECT id, room_id, username, content, created_at
		 FROM roommessages
		 WHERE room_id = $1 AND id < $3
		 ORDER BY created_at DESC
		 LIMIT $2`

const SendDirectMessageQuery = `INSERT INTO directmessages (sender,receiver,content)
								VALUES ($1,$2,$3)`

const GetRecentDirectMessagesQuery = `SELECT id,receiver,sender,content,created_at
										FROM directmessages
										WHERE sender=$1 OR receiver=$1
										ORDER BY created_at DESC
										LIMIT $2`

const GetRangedDirectMessagesQuery = `SELECT id,receiver,sender,content,created_at
										FROM directmessages
										WHERE (sender=$1 OR receiver=$1) AND id < $3
										ORDER BY created_at DESC
										LIMIT $2`

func (s *Store) SaveMessage(msg string, roomId int, userName string) error {
	_, err := s.db.Exec(
		SaveMessageQuery,
		roomId,
		userName,
		msg,
	)
	return err
}

func (s *Store) GetRecentMessages(roomId int, limit int, lastid int) ([]models.Message, error) {
	var rows *sql.Rows
	var err error
	if lastid == 0 {
		rows, err = s.db.Query(
			GetRecentMessagesQuery,
			roomId,
			limit,
		)
	} else {
		rows, err = s.db.Query(
			GetRangedMessagesQuery,
			roomId,
			limit,
			lastid,
		)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(
			&m.ID,
			&m.Room,
			&m.User,
			&m.Content,
			&m.CreatedAt,
		); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}

func (s *Store) SendDirectMessage(msg string, receiver string, user string) error {
	_, err := s.db.Exec(
		SendDirectMessageQuery, user, receiver, msg,
	)
	return err
}

func (s *Store) GetRecentDirectMessages(username string, limit int, lastid int) ([]models.Message, error) {

	var rows *sql.Rows
	var err error
	if lastid == 0 {
		rows, err = s.db.Query(
			GetRecentDirectMessagesQuery,
			username,
			limit,
		)
	} else {
		rows, err = s.db.Query(
			GetRangedDirectMessagesQuery,
			username,
			limit,
			lastid,
		)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(
			&m.ID,
			&m.Receiver,
			&m.User,
			&m.Content,
			&m.CreatedAt,
		); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil

}
