package responses

import "time"

type EventResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	UserID      string    `json:"user_id"`
	RemindFor   *uint32   `json:"remind_for"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
