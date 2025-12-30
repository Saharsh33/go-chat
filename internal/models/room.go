package models

import "time"

type StoredRoom struct {
	ID        int
	Name      string
	//Description   string
	CreatedAt time.Time
}