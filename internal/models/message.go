package models

import "time"

type Message struct {
	ID        int
	Type      string
	Room      int
	User      string
	Content   string
	CreatedAt time.Time
}
