package postgres

import (
	//"chat-server/internal/storage/postgres"
	"database/sql"
	//"chat-server/internal/models"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *Postgres) *Store {
	return &Store{db: db.DB}
}
