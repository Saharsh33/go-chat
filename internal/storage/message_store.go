package storage

import "chat-server/internal/models"

type MessageStore interface {
	SaveMessage(msg models.Message) error

	//limit=20=>last 20 messages
	GetRecentMessages(room string, limit int) ([]models.Message, error)
}
