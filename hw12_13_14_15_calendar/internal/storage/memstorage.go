package storage

import (
	"context"
	"sync"
	"time"

	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
)

type MemStorage struct {
	events map[uuid.UUID]Event
	mu     sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{events: make(map[uuid.UUID]Event)}
}

func (s *MemStorage) Connect(_ context.Context) error {
	return nil
}

func (s *MemStorage) Close() error {
	return nil
}

func (s *MemStorage) CreateEvent(_ context.Context, event Event) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; ok {
		return false, ErrUnableDuplicate
	}

	s.events[event.ID] = event

	return true, nil
}

func (s *MemStorage) UpdateEvent(_ context.Context, eventID uuid.UUID, event Event) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[eventID]; !ok {
		return false, ErrNotFound
	}

	s.events[eventID] = event

	return true, nil
}

func (s *MemStorage) DeleteEvent(_ context.Context, eventID uuid.UUID) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[eventID]; !ok {
		return false, ErrNotFound
	}

	delete(s.events, eventID)

	return true, nil
}

func (s *MemStorage) ListEventsDay(_ context.Context, userID uuid.UUID, date time.Time) map[uuid.UUID]Event {
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

func (s *MemStorage) ListEventsWeek(_ context.Context, userID uuid.UUID, date time.Time) map[uuid.UUID]Event {
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

func (s *MemStorage) ListEventsMonth(_ context.Context, userID uuid.UUID, date time.Time) map[uuid.UUID]Event {
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

func (s *MemStorage) listBy(
	userID uuid.UUID,
	date time.Time,
	cmp func(date, startAt, endAt time.Time) bool,
) map[uuid.UUID]Event {
	truncateDate := func(date time.Time) time.Time {
		utc := date.UTC()
		return time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
	}

	utc := truncateDate(date)
	result := make(map[uuid.UUID]Event)
	for _, event := range s.events {
		if event.UserID.String() != userID.String() {
			continue
		}

		if cmp(utc, truncateDate(event.StartAt), truncateDate(event.EndAt)) {
			result[event.ID] = event
		}
	}

	return result
}

func (s *MemStorage) ListToSendNotifies(_ context.Context, _ time.Time) ([]Notify, error) {
	return nil, nil
}

func (s *MemStorage) RemoveOldEvents(_ context.Context, _ time.Time) error {
	return nil
}

func (s *MemStorage) MarkInQueue(_ context.Context, _ uuid.UUID) error {
	return nil
}

func (s *MemStorage) MarkAsDispatched(_ context.Context, _ uuid.UUID) error {
	return nil
}
