package storage

import (
	"context"
	"time"

	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
)

type Storage interface {
	CreateEvent(context.Context, Event) bool
	UpdateEvent(context.Context, uuid.UUID, Event) bool
	DeleteEvent(context.Context, uuid.UUID) bool

	ListEventsDay(context.Context, uuid.UUID, time.Time) map[uuid.UUID]Event
	ListEventsWeek(context.Context, uuid.UUID, time.Time) map[uuid.UUID]Event
	ListEventsMonth(context.Context, uuid.UUID, time.Time) map[uuid.UUID]Event

	Close() error
}

func NewConnect(ctx context.Context, driver, dsn string) (Storage, error) {
	if driver == "postgres" {
		pg := NewPgStorage(dsn)
		if err := pg.Connect(ctx); err != nil {
			return nil, err
		}

		return pg, nil
	}

	return NewMemStorage(), nil
}
