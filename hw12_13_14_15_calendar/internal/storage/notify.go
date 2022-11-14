package storage

import (
	"time"

	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
)

type Notify struct {
	EventID     uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	StartAt     time.Time `db:"start_at"`
	EndAt       time.Time `db:"end_at"`
	UserID      uuid.UUID `db:"user_id"`
}
