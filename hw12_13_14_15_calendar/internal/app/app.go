package app

import (
	"context"
	"time"

	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Error(string)
	Warning(string)
	Info(string)
	Debug(string)
}

type Storage interface {
	CreateEvent(event storage.Event) bool
	UpdateEvent(uuid.UUID, storage.Event) bool
	DeleteEvent(uuid.UUID) bool

	ListEventsDay(uuid.UUID, time.Time) map[uuid.UUID]storage.Event
	ListEventsWeek(uuid.UUID, time.Time) map[uuid.UUID]storage.Event
	ListEventsMonth(uuid.UUID, time.Time) map[uuid.UUID]storage.Event
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
