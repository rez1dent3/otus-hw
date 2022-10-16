package storage_test

import (
	"testing"
	"time"

	storage2 "github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
	"github.com/stretchr/testify/require"
)

func TestMemStorage_Create(t *testing.T) {
	t.Run("CreateEvent.notExists", func(t *testing.T) {
		storage := storage2.NewMemStorage()
		event := storage2.Event{}

		require.Len(t, storage.ListEventsMonth(event.UserID, event.StartAt), 0)
		success := storage.CreateEvent(storage2.Event{})

		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(event.UserID, event.StartAt), 1)
	})

	t.Run("CreateEvent.exists", func(t *testing.T) {
		storage := storage2.NewMemStorage()
		event := storage2.Event{}

		require.Len(t, storage.ListEventsMonth(event.UserID, event.StartAt), 0)
		success := storage.CreateEvent(storage2.Event{})

		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(event.UserID, event.StartAt), 1)

		loss := storage.CreateEvent(storage2.Event{})
		require.False(t, loss)
		require.Len(t, storage.ListEventsMonth(event.UserID, event.StartAt), 1)
	})
}

func TestMemStorage_Update(t *testing.T) {
	t.Run("UpdateEvent", func(t *testing.T) {
		storage := storage2.NewMemStorage()
		event1 := storage2.Event{}
		loss := storage.UpdateEvent(event1.ID, event1)
		require.False(t, loss)

		success := storage.CreateEvent(event1)
		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(event1.UserID, event1.StartAt), 1)

		event2 := storage2.Event{ID: event1.ID, Title: "hello"}
		success = storage.UpdateEvent(event2.ID, event2)
		require.Len(t, storage.ListEventsMonth(event2.UserID, event2.StartAt), 1)
		require.True(t, success)
	})
}

func TestMemStorage_Delete(t *testing.T) {
	t.Run("DeleteEvent.notExists", func(t *testing.T) {
		storage := storage2.NewMemStorage()
		event1 := storage2.Event{}
		require.Len(t, storage.ListEventsMonth(event1.UserID, event1.StartAt), 0)
		success := storage.DeleteEvent(event1.ID)
		require.False(t, success)
	})

	t.Run("DeleteEvent.exists", func(t *testing.T) {
		storage := storage2.NewMemStorage()
		event1 := storage2.Event{}
		success := storage.CreateEvent(event1)
		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(event1.UserID, event1.StartAt), 1)

		success = storage.DeleteEvent(event1.ID)
		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(event1.UserID, event1.StartAt), 0)
	})
}

func TestMemStorage_List(t *testing.T) {
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

		storage := storage2.NewMemStorage()
		for _, input := range inputs {
			storage.CreateEvent(
				storage2.Event{ID: uuid.FromString(input.eventID),
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

		storage := storage2.NewMemStorage()
		for _, input := range inputs {
			storage.CreateEvent(storage2.Event{ID: uuid.FromString(input.uuid), UserID: userID, StartAt: input.date, EndAt: input.date})
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

		storage := storage2.NewMemStorage()
		for _, input := range inputs {
			storage.CreateEvent(storage2.Event{ID: uuid.FromString(input.uuid), UserID: userID, StartAt: input.date, EndAt: input.date})
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

		storage := storage2.NewMemStorage()
		for _, input := range inputs {
			storage.CreateEvent(storage2.Event{ID: uuid.FromString(input.uuid), UserID: userID, StartAt: input.date, EndAt: input.date})
		}

		require.Len(t, storage.ListEventsMonth(userID, now.AddDate(0, 1, 0)), 1)
		require.Len(t, storage.ListEventsMonth(userID, now), 1)
		require.Len(t, storage.ListEventsMonth(userID, now.AddDate(0, 1, 0)), 1)
	})
}
