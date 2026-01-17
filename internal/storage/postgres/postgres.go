package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
}

func NewDB(dsn string) *Postgres {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("db open error:", err)
		return nil
	}

	if err := db.Ping(); err != nil {
		log.Fatal("DB Ping error:- ", err)
		return nil
	}

	return &Postgres{DB: db}
}
