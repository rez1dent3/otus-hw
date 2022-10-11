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
	return &Storage{}
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

	delete(s.events, eventId)

	return true
}

func (s *Storage) ListEventToday(date time.Time) []storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return []storage.Event{}
}

func (s *Storage) ListEventWeek(date time.Time) []storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return []storage.Event{}
}

func (s *Storage) ListEventMonth(date time.Time) []storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return []storage.Event{}
}
