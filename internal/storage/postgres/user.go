package postgres

import (
	"context"
	"log"
)

const (
	CreateUserIfNotExistsQuery = `INSERT INTO users (username) VALUES ($1) ON CONFLICT (username) DO NOTHING;`
	
	GetUserByNameQuery = `SELECT Id FROM users WHERE username = $1`
)

func (s *Store) CreateUserIfNotExists(ctx context.Context, user string) {

	_, err := s.db.Exec(CreateUserIfNotExistsQuery, user)

	if err != nil {
		log.Println(err)
		return
	}

}

func (s *Store) GetUserByName(ctx context.Context, user string) (int, error) {
	userId, err := s.db.QueryContext(ctx, GetUserByNameQuery, user)
	var id int
	for userId.Next() {
		if err := userId.Scan(
			&id,
		); err != nil {
			return 0, err
		}
	}
	return id, err
}
