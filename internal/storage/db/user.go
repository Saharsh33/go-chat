package postgres

import (
	"context"
	"log"
)

const (
	CreateUserIfNotExistsQuery = `INSERT INTO users (username) VALUES ($1) ON CONFLICT (username) DO NOTHING;`

	GetUserByNameQuery = `SELECT Id FROM users WHERE username = $1`

	RegisterUserQuery = `INSERT INTO users (username, password_hash) VALUES ($1, $2)
		ON CONFLICT (username) DO NOTHING`

	GetUserPasswordHashQuery = `SELECT COALESCE(password_hash, '') FROM users WHERE username = $1`
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

// RegisterUser inserts a new user with the given bcrypt password hash.
// Returns an error if the username is already taken.
func (s *Store) RegisterUser(ctx context.Context, username, passwordHash string) error {
	result, err := s.db.ExecContext(ctx, RegisterUserQuery, username, passwordHash)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserExists
	}
	return nil
}

// GetUserPasswordHash returns the bcrypt hash stored for the given username.
// Returns an error if the user does not exist.
func (s *Store) GetUserPasswordHash(ctx context.Context, username string) (string, error) {
	var hash string
	err := s.db.QueryRowContext(ctx, GetUserPasswordHashQuery, username).Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}

