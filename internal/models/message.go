package models

import "time"

type Message struct {
	ID        int
	Type      string
	Room      string
	User      string
	Content   string
	CreatedAt time.Time
}
