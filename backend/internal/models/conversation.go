package models

import "time"

type Conversation struct {
	ID        int
	UserOne   string
	UserTwo   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
