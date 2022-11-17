package storage

import (
	"context"
	"errors"
	"time"

	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
)

var (
	ErrUnableDuplicate = errors.New("unable to duplicate")
	ErrNotFound        = errors.New("not found")
)

type CalendarStorage interface {
	CreateEvent(context.Context, Event) (bool, error)
	UpdateEvent(context.Context, uuid.UUID, Event) (bool, error)
	DeleteEvent(context.Context, uuid.UUID) (bool, error)

	ListEventsDay(context.Context, uuid.UUID, time.Time) map[uuid.UUID]Event
	ListEventsWeek(context.Context, uuid.UUID, time.Time) map[uuid.UUID]Event
	ListEventsMonth(context.Context, uuid.UUID, time.Time) map[uuid.UUID]Event

	Close() error
}

type SchedulerStorage interface {
	ListToSendNotifies(context.Context, time.Time) ([]Notify, error)
	RemoveOldEvents(context.Context, time.Time) error
	MarkInQueue(context.Context, uuid.UUID) error

	Close() error
}

type SenderStorage interface {
	MarkAsDispatched(context.Context, uuid.UUID) error

	Close() error
}

type internalStorage interface {
	CalendarStorage
	SchedulerStorage
	SenderStorage
}

func newConnect(ctx context.Context, driver, dsn string) (internalStorage, error) {
	if driver == "postgres" {
		pg := NewPgStorage(dsn)
		if err := pg.Connect(ctx); err != nil {
			return nil, err
		}

		return pg, nil
	}

	return NewMemStorage(), nil
}

func NewCalendarStorage(ctx context.Context, driver, dsn string) (CalendarStorage, error) {
	return newConnect(ctx, driver, dsn)
}

func NewSchedulerStorage(ctx context.Context, driver, dsn string) (SchedulerStorage, error) {
	return newConnect(ctx, driver, dsn)
}

func NewSenderStorage(ctx context.Context, driver, dsn string) (SenderStorage, error) {
	return newConnect(ctx, driver, dsn)
}
