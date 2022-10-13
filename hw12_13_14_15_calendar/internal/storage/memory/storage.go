package memorystorage

import (
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
	"sync"
	"time"
)

type Storage struct {
	events map[uuid.UUID]storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{events: make(map[uuid.UUID]storage.Event)}
}

func (s *Storage) CreateEvent(event storage.Event) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[event.ID] = event

	return true
}

func (s *Storage) UpdateEvent(eventId uuid.UUID, event storage.Event) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[eventId]; !ok {
		return false
	}

	s.events[eventId] = event

	return true
}

func (s *Storage) DeleteEvent(eventId uuid.UUID) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[eventId]; !ok {
		return false
	}

	delete(s.events, eventId)

	return true
}

func (s *Storage) ListEventsDay(date time.Time) map[uuid.UUID]storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.listBy(
		date,
		func(date, startAt, endAt time.Time) bool {
			if date.Equal(startAt) || date.Equal(endAt) {
				return true
			}

			return date.After(startAt) && date.Before(endAt)
		},
	)
}

func (s *Storage) ListEventsWeek(date time.Time) map[uuid.UUID]storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.listBy(
		date,
		func(date, startAt, endAt time.Time) bool {
			startYear, startWeek := startAt.ISOWeek()
			endYear, endWeek := endAt.ISOWeek()
			dateYear, dateWeek := date.ISOWeek()
			week := weekStart(dateYear, dateWeek)
			start := weekStart(startYear, startWeek)
			end := weekStart(endYear, endWeek)

			if week.Equal(start) || week.Equal(end) {
				return true
			}

			return week.After(start) && week.Before(end)
		},
	)
}

func (s *Storage) ListEventsMonth(date time.Time) map[uuid.UUID]storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.listBy(
		date,
		func(date, startAt, endAt time.Time) bool {
			if date.Equal(startAt) || date.Equal(endAt) {
				return true
			}

			return date.After(startAt) && date.Before(endAt)
		},
	)
}

func weekStart(year, week int) time.Time {
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}

func (s *Storage) listBy(date time.Time, cmp func(date, startAt, endAt time.Time) bool) map[uuid.UUID]storage.Event {
	truncateDate := func(date time.Time) time.Time {
		utc := date.UTC()
		return time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
	}

	utc := truncateDate(date)
	result := make(map[uuid.UUID]storage.Event)
	for _, event := range s.events {
		if cmp(utc, truncateDate(event.StartAt), truncateDate(event.EndAt)) {
			result[event.ID] = event
		}
	}

	return result
}
