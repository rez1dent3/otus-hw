package storage

import (
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
	"sync"
	"time"
)

type MemStorage struct {
	events map[uuid.UUID]Event
	mu     sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{events: make(map[uuid.UUID]Event)}
}

func (s *MemStorage) CreateEvent(event Event) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; ok {
		return false
	}

	s.events[event.ID] = event

	return true
}

func (s *MemStorage) UpdateEvent(eventID uuid.UUID, event Event) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[eventID]; !ok {
		return false
	}

	s.events[eventID] = event

	return true
}

func (s *MemStorage) DeleteEvent(eventID uuid.UUID) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[eventID]; !ok {
		return false
	}

	delete(s.events, eventID)

	return true
}

func (s *MemStorage) ListEventsDay(userID uuid.UUID, date time.Time) map[uuid.UUID]Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.listBy(
		userID,
		date,
		func(date, startAt, endAt time.Time) bool {
			if date.Equal(startAt) || date.Equal(endAt) {
				return true
			}

			return date.After(startAt) && date.Before(endAt)
		},
	)
}

func (s *MemStorage) ListEventsWeek(userID uuid.UUID, date time.Time) map[uuid.UUID]Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.listBy(
		userID,
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

func (s *MemStorage) ListEventsMonth(userID uuid.UUID, date time.Time) map[uuid.UUID]Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.listBy(
		userID,
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

func (s *MemStorage) listBy(userID uuid.UUID, date time.Time, cmp func(date, startAt, endAt time.Time) bool) map[uuid.UUID]Event {
	truncateDate := func(date time.Time) time.Time {
		utc := date.UTC()
		return time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
	}

	utc := truncateDate(date)
	result := make(map[uuid.UUID]Event)
	for _, event := range s.events {
		if event.UserID.ToString() != userID.ToString() {
			continue
		}

		if cmp(utc, truncateDate(event.StartAt), truncateDate(event.EndAt)) {
			result[event.ID] = event
		}
	}

	return result
}
