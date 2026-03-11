package postgres

import (
	//"chat-server/internal/storage/postgres"
	"database/sql"
	"errors"
	//"chat-server/internal/models"
)

// ErrUserExists is returned when attempting to register a username that is already taken.
var ErrUserExists = errors.New("username already taken")

type Store struct {
	db *sql.DB
}

func NewStore(db *Postgres) *Store {
	return &Store{db: db.DB}
}
