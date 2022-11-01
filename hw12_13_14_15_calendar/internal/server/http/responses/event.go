package responses

import "time"

type EventResponse struct {
	ID          string
	Title       string
	Description *string
	StartAt     time.Time
	EndAt       time.Time
	UserID      string
	RemindFor   *uint32
}

type ErrorResponse struct {
	Message string
}
