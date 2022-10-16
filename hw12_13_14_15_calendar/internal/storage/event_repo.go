package storage

import (
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
	"time"
)

type EventRepoInterface interface {
	CreateEvent(Event) bool
	UpdateEvent(uuid.UUID, Event) bool
	DeleteEvent(uuid.UUID) bool

	ListEventsDay(uuid.UUID, time.Time) map[uuid.UUID]Event
	ListEventsWeek(uuid.UUID, time.Time) map[uuid.UUID]Event
	ListEventsMonth(uuid.UUID, time.Time) map[uuid.UUID]Event
}
