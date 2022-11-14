package requests

import (
	"time"
)

type EventRequest struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	StartAt     time.Time      `json:"start_at"`
	EndAt       time.Time      `json:"end_at"`
	UserID      string         `json:"user_id"`
	RemindFor   *time.Duration `json:"remind_for"`
}

type EventListRequest struct {
	UserID string    `json:"user_id"`
	Date   time.Time `json:"date"`
}
