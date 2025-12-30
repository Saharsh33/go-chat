package postgres

import (
	"chat-server/internal/models"
)


const SaveMessageQuery = `INSERT INTO messages (room, username, content)
		 VALUES ($1, $2, $3)`

const GetRecentMessagesQuery = `SELECT id, room, username, content, created_at
		 FROM messages
		 WHERE room = $1
		 ORDER BY created_at DESC
		 LIMIT $2`

		 
func (s *Store) SaveMessage(msg models.Message) error {
	_, err := s.db.Exec(
		SaveMessageQuery,
		msg.Room,
		msg.User,
		msg.Content,
	)
	return err
}

func (s *Store) GetRecentMessages(room string, limit int) ([]models.Message, error) {
	rows, err := s.db.Query(
		GetRecentMessagesQuery,
		room,
		limit,
	)
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
