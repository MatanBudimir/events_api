package models

import (
	"time"
)

type EventMeeting struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	User        User      `json:"-"`
	Event       Event     `json:"-"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Scheduled   bool      `json:"scheduled"`
}
