package models

import "time"

type RoomMember struct {
	RoomID   int
	Username string
	JoinedAt time.Time
}