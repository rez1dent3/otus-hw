package app

import (
	"context"
	"errors"
	"time"

	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
)

var (
	ErrCreateEvent = errors.New("create event error")
	ErrUpdateEvent = errors.New("update event error")
)

type Application interface {
	CreateEvent(
		ctx context.Context,
		id string,
		title string,
		description string,
		startAt time.Time,
		endAt time.Time,
		userID string,
		remindFor uint32,
	) error

	UpdateEvent(
		ctx context.Context,
		id string,
		title string,
		description string,
		startAt time.Time,
		endAt time.Time,
		userID string,
		remindFor uint32,
	) error

	DeleteEvent(ctx context.Context, id string) bool

	ListEventsDay(context.Context, string, time.Time) []storage.Event
	ListEventsWeek(context.Context, string, time.Time) []storage.Event
	ListEventsMonth(context.Context, string, time.Time) []storage.Event
}

type App struct {
	logger  *logger.Logger
	storage Storage
}

type Storage interface {
	CreateEvent(context.Context, storage.Event) bool
	UpdateEvent(context.Context, uuid.UUID, storage.Event) bool
	DeleteEvent(context.Context, uuid.UUID) bool

	ListEventsDay(context.Context, uuid.UUID, time.Time) map[uuid.UUID]storage.Event
	ListEventsWeek(context.Context, uuid.UUID, time.Time) map[uuid.UUID]storage.Event
	ListEventsMonth(context.Context, uuid.UUID, time.Time) map[uuid.UUID]storage.Event
}

func New(logger *logger.Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) CreateEvent(
	ctx context.Context,
	id string,
	title string,
	description string,
	startAt time.Time,
	endAt time.Time,
	userID string,
	remindFor uint32,
) error {
	result := a.storage.CreateEvent(ctx, storage.Event{
		ID:          uuid.FromString(id),
		Title:       title,
		Description: description,
		StartAt:     startAt,
		EndAt:       endAt,
		UserID:      uuid.FromString(userID),
		RemindFor:   &remindFor,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	if !result {
		return ErrCreateEvent
	}

	return nil
}

func (a *App) UpdateEvent(
	ctx context.Context,
	id string,
	title string,
	description string,
	startAt time.Time,
	endAt time.Time,
	userID string,
	remindFor uint32,
) error {
	result := a.storage.UpdateEvent(ctx, uuid.FromString(id), storage.Event{
		Title:       title,
		Description: description,
		StartAt:     startAt,
		EndAt:       endAt,
		UserID:      uuid.FromString(userID),
		RemindFor:   &remindFor,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	if !result {
		return ErrUpdateEvent
	}

	return nil
}

func (a *App) DeleteEvent(
	ctx context.Context,
	id string,
) bool {
	return a.storage.DeleteEvent(ctx, uuid.FromString(id))
}

func (a *App) ListEventsDay(ctx context.Context, userID string, date time.Time) []storage.Event {
	return a.listResponse(a.storage.ListEventsDay(ctx, uuid.FromString(userID), date))
}

func (a *App) ListEventsWeek(ctx context.Context, userID string, date time.Time) []storage.Event {
	return a.listResponse(a.storage.ListEventsWeek(ctx, uuid.FromString(userID), date))
}

func (a *App) ListEventsMonth(ctx context.Context, userID string, date time.Time) []storage.Event {
	return a.listResponse(a.storage.ListEventsMonth(ctx, uuid.FromString(userID), date))
}

func (a *App) listResponse(events map[uuid.UUID]storage.Event) []storage.Event {
	i := 0
	result := make([]storage.Event, len(events))
	for _, item := range events {
		result[i] = item
		i++
	}

	return result
}
