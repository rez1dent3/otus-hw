package requests

import (
	"time"
)

type EventRequest struct {
	ID          string
	Title       string
	Description string
	StartAt     time.Time
	EndAt       time.Time
	UserID      string
	RemindFor   uint32
}

type EventListRequest struct {
	UserID string
	Date   time.Time
}
