package storage

import (
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
	"time"
)

type Event struct {
	ID          uuid.UUID
	Title       string
	Description string
	StartAt     time.Time
	EndAt       time.Time
	UserID      uuid.UUID
	RemindFor   uint32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
