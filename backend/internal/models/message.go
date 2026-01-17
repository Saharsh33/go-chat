package models

import "time"

type Message struct {
	ID        int
	Type      string
	Room      int
	Receiver  string
	User      string
	Content   string
	CreatedAt time.Time
}
