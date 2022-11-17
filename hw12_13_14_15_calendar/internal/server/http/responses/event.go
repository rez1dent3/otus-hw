package responses

import "time"

type EventResponse struct {
	ID           string         `json:"id"`
	Title        string         `json:"title"`
	Description  *string        `json:"description"`
	StartAt      time.Time      `json:"start_at"`
	EndAt        time.Time      `json:"end_at"`
	UserID       string         `json:"user_id"`
	RemindFor    *time.Duration `json:"remind_for"`
	IsDispatched bool           `json:"is_dispatched"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
