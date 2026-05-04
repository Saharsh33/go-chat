package postgres

import (
	"chat-server/internal/models"
	"context"
	"errors"
	"log"
)

const (
	// GetOrCreateConversationQuery upserts a conversation for the canonical pair.
	// user_one is always lexicographically smaller than user_two.
	GetOrCreateConversationQuery = `INSERT INTO conversations (user_one, user_two)
		VALUES (LEAST($1,$2), GREATEST($1,$2))
		ON CONFLICT (user_one, user_two) DO UPDATE SET updated_at = NOW()
		RETURNING id, user_one, user_two, created_at, updated_at`

	GetConversationsOfUserQuery = `SELECT id, user_one, user_two, created_at, updated_at
		FROM conversations
		WHERE user_one = $1 OR user_two = $1
		ORDER BY updated_at DESC`
)

func (s *Store) GetOrCreateConversation(ctx context.Context, userA string, userB string) (*models.Conversation, error) {
	if userA == userB {
		return nil, errors.New("cannot create conversation with yourself")
	}
	var c models.Conversation
	err := s.db.QueryRowContext(ctx, GetOrCreateConversationQuery, userA, userB).Scan(
		&c.ID, &c.UserOne, &c.UserTwo, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		log.Println("GetOrCreateConversation error:", err)
		return nil, err
	}
	return &c, nil
}

func (s *Store) GetConversationsOfUser(ctx context.Context, username string) ([]models.Conversation, error) {
	rows, err := s.db.QueryContext(ctx, GetConversationsOfUserQuery, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convos []models.Conversation
	for rows.Next() {
		var c models.Conversation
		if err := rows.Scan(&c.ID, &c.UserOne, &c.UserTwo, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		convos = append(convos, c)
	}
	return convos, nil
}
