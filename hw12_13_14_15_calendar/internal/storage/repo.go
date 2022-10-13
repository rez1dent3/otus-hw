package storage

import (
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
	"time"
)

type EventRepoInterface interface {
	CreateEvent(event Event) bool
	UpdateEvent(eventId uuid.UUID, event Event) bool
	DeleteEvent(eventId uuid.UUID) bool

	ListEventsDay(date time.Time) map[uuid.UUID]Event
	ListEventsWeek(date time.Time) map[uuid.UUID]Event
	ListEventsMonth(date time.Time) map[uuid.UUID]Event
}
