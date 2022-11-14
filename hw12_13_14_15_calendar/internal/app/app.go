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
		remindFor *time.Duration,
	) error

	UpdateEvent(
		ctx context.Context,
		id string,
		title string,
		description string,
		startAt time.Time,
		endAt time.Time,
		userID string,
		remindFor *time.Duration,
	) error

	DeleteEvent(ctx context.Context, id string) error

	ListEventsDay(context.Context, string, time.Time) []storage.Event
	ListEventsWeek(context.Context, string, time.Time) []storage.Event
	ListEventsMonth(context.Context, string, time.Time) []storage.Event
}

type App struct {
	logger  *logger.Logger
	storage storage.CalendarStorage
}

func New(logger *logger.Logger, storage storage.CalendarStorage) *App {
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
	remindFor *time.Duration,
) error {
	if result, err := a.storage.CreateEvent(ctx, storage.Event{
		ID:          uuid.FromString(id),
		Title:       title,
		Description: description,
		StartAt:     startAt,
		EndAt:       endAt,
		UserID:      uuid.FromString(userID),
		RemindFor:   (*storage.Duration)(remindFor),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}); err != nil {
		return err
	} else if !result {
		return ErrCreateEvent
	} else {
		return nil
	}
}

func (a *App) UpdateEvent(
	ctx context.Context,
	id string,
	title string,
	description string,
	startAt time.Time,
	endAt time.Time,
	userID string,
	remindFor *time.Duration,
) error {
	idObj := uuid.FromString(id)
	if result, err := a.storage.UpdateEvent(ctx, idObj, storage.Event{
		ID:          idObj,
		Title:       title,
		Description: description,
		StartAt:     startAt,
		EndAt:       endAt,
		UserID:      uuid.FromString(userID),
		RemindFor:   (*storage.Duration)(remindFor),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}); err != nil {
		return err
	} else if !result {
		return ErrUpdateEvent
	} else {
		return nil
	}
}

func (a *App) DeleteEvent(
	ctx context.Context,
	id string,
) error {
	if result, err := a.storage.DeleteEvent(ctx, uuid.FromString(id)); err != nil {
		return err
	} else if !result {
		return ErrUpdateEvent
	} else {
		return nil
	}
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
