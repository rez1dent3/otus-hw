package app

import (
	"context"
	"time"

	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
)

type Application interface {
	CreateEvent(ctx context.Context, id, title string) error
}

type App struct {
	logger  *logger.Logger
	storage Storage
}

type Storage interface {
	CreateEvent(event storage.Event) bool
	UpdateEvent(uuid.UUID, storage.Event) bool
	DeleteEvent(uuid.UUID) bool

	ListEventsDay(uuid.UUID, time.Time) map[uuid.UUID]storage.Event
	ListEventsWeek(uuid.UUID, time.Time) map[uuid.UUID]storage.Event
	ListEventsMonth(uuid.UUID, time.Time) map[uuid.UUID]storage.Event
}

func New(logger *logger.Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
