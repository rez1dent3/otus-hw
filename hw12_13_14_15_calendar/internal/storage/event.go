package storage

import (
	"time"

	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
)

type Event struct {
	ID           uuid.UUID `db:"id"`
	Title        string    `db:"title"`
	Description  string    `db:"description"`
	StartAt      time.Time `db:"start_at"`
	EndAt        time.Time `db:"end_at"`
	UserID       uuid.UUID `db:"user_id"`
	RemindFor    *Duration `db:"remind_for"`
	InQueue      bool      `db:"in_queue"`
	IsDispatched bool      `db:"is_dispatched"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
