package storage

import (
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestStorage_CRUD(t *testing.T) {
	t.Run("CreateEvent.notExists", func(t *testing.T) {
		storage := NewMemStorage()
		event := Event{}

		require.Len(t, storage.ListEventsMonth(event.UserID, event.StartAt), 0)
		success := storage.CreateEvent(Event{})

		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(event.UserID, event.StartAt), 1)
	})

	t.Run("CreateEvent.exists", func(t *testing.T) {
		storage := NewMemStorage()
		event := Event{}

		require.Len(t, storage.ListEventsMonth(event.UserID, event.StartAt), 0)
		success := storage.CreateEvent(Event{})

		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(event.UserID, event.StartAt), 1)

		loss := storage.CreateEvent(Event{})
		require.False(t, loss)
		require.Len(t, storage.ListEventsMonth(event.UserID, event.StartAt), 1)
	})

	t.Run("UpdateEvent", func(t *testing.T) {
		storage := NewMemStorage()
		event1 := Event{}
		loss := storage.UpdateEvent(event1.ID, event1)
		require.False(t, loss)

		success := storage.CreateEvent(event1)
		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(event1.UserID, event1.StartAt), 1)

		event2 := Event{ID: event1.ID, Title: "hello"}
		success = storage.UpdateEvent(event2.ID, event2)
		require.Len(t, storage.ListEventsMonth(event2.UserID, event2.StartAt), 1)
		require.True(t, success)
	})

	t.Run("DeleteEvent.notExists", func(t *testing.T) {
		storage := NewMemStorage()
		event1 := Event{}
		require.Len(t, storage.ListEventsMonth(event1.UserID, event1.StartAt), 0)
		success := storage.DeleteEvent(event1.ID)
		require.False(t, success)
	})

	t.Run("DeleteEvent.exists", func(t *testing.T) {
		storage := NewMemStorage()
		event1 := Event{}
		success := storage.CreateEvent(event1)
		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(event1.UserID, event1.StartAt), 1)

		success = storage.DeleteEvent(event1.ID)
		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(event1.UserID, event1.StartAt), 0)
	})

	t.Run("list.filterByUserID", func(t *testing.T) {
		now := time.Now()
		inputs := []struct {
			eventID string
			userID  string
			date    time.Time
		}{
			{eventID: "85715fcaeb2d470a90b044f2874933a4", userID: "571548163f594e5f9fcaba2661a1ac88", date: now.AddDate(0, 0, -1)},
			{eventID: "87a74e3a145840f4aa6874b7c1bfa779", userID: "571548163f594e5f9fcaba2661a1ac88", date: now},
			{eventID: "6f2742e521944580af42533708ff2970", userID: "69d3583e60c2436b85679fb836e6aab2", date: now.AddDate(0, 0, 1)},
		}

		storage := NewMemStorage()
		for _, input := range inputs {
			storage.CreateEvent(
				Event{ID: uuid.FromString(input.eventID),
					UserID:  uuid.FromString(input.userID),
					StartAt: input.date,
					EndAt:   input.date})
		}

		require.Len(t, storage.ListEventsDay(uuid.FromString("571548163f594e5f9fcaba2661a1ac88"), now.AddDate(0, 0, -1)), 1)
		require.Len(t, storage.ListEventsDay(uuid.FromString("571548163f594e5f9fcaba2661a1ac88"), now), 1)
		require.Len(t, storage.ListEventsDay(uuid.FromString("571548163f594e5f9fcaba2661a1ac88"), now.AddDate(0, 0, 1)), 0)

		require.Len(t, storage.ListEventsDay(uuid.FromString("69d3583e60c2436b85679fb836e6aab2"), now.AddDate(0, 0, -1)), 0)
		require.Len(t, storage.ListEventsDay(uuid.FromString("69d3583e60c2436b85679fb836e6aab2"), now), 0)
		require.Len(t, storage.ListEventsDay(uuid.FromString("69d3583e60c2436b85679fb836e6aab2"), now.AddDate(0, 0, 1)), 1)
	})

	t.Run("ListEventsDay", func(t *testing.T) {
		now := time.Now()
		inputs := []struct {
			uuid string
			date time.Time
		}{
			{uuid: "85715fcaeb2d470a90b044f2874933a4", date: now.AddDate(0, 0, -1)},
			{uuid: "87a74e3a145840f4aa6874b7c1bfa779", date: now},
			{uuid: "6f2742e521944580af42533708ff2970", date: now.AddDate(0, 0, 1)},
		}

		userID := uuid.FromString("571548163f594e5f9fcaba2661a1ac88")

		storage := NewMemStorage()
		for _, input := range inputs {
			storage.CreateEvent(Event{ID: uuid.FromString(input.uuid), UserID: userID, StartAt: input.date, EndAt: input.date})
		}

		require.Len(t, storage.ListEventsDay(userID, now.AddDate(0, 0, -1)), 1)
		require.Len(t, storage.ListEventsDay(userID, now), 1)
		require.Len(t, storage.ListEventsDay(userID, now.AddDate(0, 0, 1)), 1)
	})

	t.Run("ListEventsWeek", func(t *testing.T) {
		now := time.Now()
		inputs := []struct {
			uuid string
			date time.Time
		}{
			{uuid: "85715fcaeb2d470a90b044f2874933a4", date: now.AddDate(0, 0, -7)},
			{uuid: "87a74e3a145840f4aa6874b7c1bfa779", date: now},
			{uuid: "6f2742e521944580af42533708ff2970", date: now.AddDate(0, 0, 7)},
		}

		userID := uuid.FromString("571548163f594e5f9fcaba2661a1ac88")

		storage := NewMemStorage()
		for _, input := range inputs {
			storage.CreateEvent(Event{ID: uuid.FromString(input.uuid), UserID: userID, StartAt: input.date, EndAt: input.date})
		}

		require.Len(t, storage.ListEventsWeek(userID, now.AddDate(0, 0, -7)), 1)
		require.Len(t, storage.ListEventsWeek(userID, now), 1)
		require.Len(t, storage.ListEventsWeek(userID, now.AddDate(0, 0, 7)), 1)
	})

	t.Run("ListEventsMonth", func(t *testing.T) {
		now := time.Now()
		inputs := []struct {
			uuid string
			date time.Time
		}{
			{uuid: "85715fcaeb2d470a90b044f2874933a4", date: now.AddDate(0, -1, 0)},
			{uuid: "87a74e3a145840f4aa6874b7c1bfa779", date: now},
			{uuid: "6f2742e521944580af42533708ff2970", date: now.AddDate(0, 1, 0)},
		}

		userID := uuid.FromString("571548163f594e5f9fcaba2661a1ac88")

		storage := NewMemStorage()
		for _, input := range inputs {
			storage.CreateEvent(Event{ID: uuid.FromString(input.uuid), UserID: userID, StartAt: input.date, EndAt: input.date})
		}

		require.Len(t, storage.ListEventsMonth(userID, now.AddDate(0, 1, 0)), 1)
		require.Len(t, storage.ListEventsMonth(userID, now), 1)
		require.Len(t, storage.ListEventsMonth(userID, now.AddDate(0, 1, 0)), 1)
	})

}