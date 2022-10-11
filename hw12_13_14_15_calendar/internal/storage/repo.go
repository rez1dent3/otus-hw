package storage

import (
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
	"time"
)

type EventRepoInterface interface {
	CreateEvent(event Event) bool
	UpdateEvent(eventId uuid.UUID, event Event) bool
	DeleteEvent(eventId uuid.UUID) bool

	ListEventToday(date time.Time) []Event
	ListEventWeek(date time.Time) []Event
	ListEventMonth(date time.Time) []Event
}
