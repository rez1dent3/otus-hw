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
		event := storage2.Event{ID: uuid.Gen(), UserID: uuid.Gen()}

		require.Len(t, storage.ListEventsMonth(event.UserID, event.StartAt), 0)
		success := storage.CreateEvent(event)

		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(event.UserID, event.StartAt), 1)
	})

	t.Run("CreateEvent.exists", func(t *testing.T) {
		storage := storage2.NewMemStorage()
		event := storage2.Event{ID: uuid.Gen(), UserID: uuid.Gen()}

		require.Len(t, storage.ListEventsMonth(event.UserID, event.StartAt), 0)
		success := storage.CreateEvent(event)

		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(event.UserID, event.StartAt), 1)

		loss := storage.CreateEvent(event)
		require.False(t, loss)
		require.Len(t, storage.ListEventsMonth(event.UserID, event.StartAt), 1)
	})
}

func TestMemStorage_Update(t *testing.T) {
	t.Run("UpdateEvent", func(t *testing.T) {
		storage := storage2.NewMemStorage()
		event1 := storage2.Event{ID: uuid.Gen(), UserID: uuid.Gen()}
		loss := storage.UpdateEvent(event1.ID, event1)
		require.False(t, loss)

		success := storage.CreateEvent(event1)
		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(event1.UserID, event1.StartAt), 1)

		event2 := storage2.Event{ID: event1.ID, UserID: uuid.Gen(), Title: "hello"}
		success = storage.UpdateEvent(event2.ID, event2)
		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(event2.UserID, event2.StartAt), 1)
	})
}

func TestMemStorage_Delete(t *testing.T) {
	t.Run("DeleteEvent.notExists", func(t *testing.T) {
		storage := storage2.NewMemStorage()
		event1 := storage2.Event{ID: uuid.Gen()}
		require.Len(t, storage.ListEventsMonth(event1.UserID, event1.StartAt), 0)
		success := storage.DeleteEvent(event1.ID)
		require.False(t, success)
	})

	t.Run("DeleteEvent.exists", func(t *testing.T) {
		storage := storage2.NewMemStorage()
		event1 := storage2.Event{ID: uuid.Gen()}
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
		storage := storage2.NewMemStorage()
		user1 := uuid.Gen()
		user2 := uuid.Gen()

		now := time.Now()
		inputs := []struct {
			eventID uuid.UUID
			userID  uuid.UUID
			date    time.Time
		}{
			{eventID: uuid.Gen(), userID: user1, date: now.AddDate(0, 0, -1)},
			{eventID: uuid.Gen(), userID: user1, date: now},
			{eventID: uuid.Gen(), userID: user2, date: now.AddDate(0, 0, 1)},
		}

		for _, input := range inputs {
			storage.CreateEvent(
				storage2.Event{ID: input.eventID,
					UserID:  input.userID,
					StartAt: input.date,
					EndAt:   input.date})
		}

		require.Len(t, storage.ListEventsDay(user1, now.AddDate(0, 0, -1)), 1)
		require.Len(t, storage.ListEventsDay(user1, now), 1)
		require.Len(t, storage.ListEventsDay(user1, now.AddDate(0, 0, 1)), 0)

		require.Len(t, storage.ListEventsDay(user2, now.AddDate(0, 0, -1)), 0)
		require.Len(t, storage.ListEventsDay(user2, now), 0)
		require.Len(t, storage.ListEventsDay(user2, now.AddDate(0, 0, 1)), 1)
	})

	t.Run("ListEventsDay", func(t *testing.T) {
		storage := storage2.NewMemStorage()
		now := time.Now()
		inputs := []struct {
			uuid uuid.UUID
			date time.Time
		}{
			{uuid: uuid.Gen(), date: now.AddDate(0, 0, -1)},
			{uuid: uuid.Gen(), date: now},
			{uuid: uuid.Gen(), date: now.AddDate(0, 0, 1)},
		}

		userID := uuid.Gen()

		for _, input := range inputs {
			storage.CreateEvent(storage2.Event{ID: input.uuid, UserID: userID, StartAt: input.date, EndAt: input.date})
		}

		require.Len(t, storage.ListEventsDay(userID, now.AddDate(0, 0, -1)), 1)
		require.Len(t, storage.ListEventsDay(userID, now), 1)
		require.Len(t, storage.ListEventsDay(userID, now.AddDate(0, 0, 1)), 1)
	})

	t.Run("ListEventsWeek", func(t *testing.T) {
		storage := storage2.NewMemStorage()
		now := time.Now()
		inputs := []struct {
			uuid uuid.UUID
			date time.Time
		}{
			{uuid: uuid.Gen(), date: now.AddDate(0, 0, -7)},
			{uuid: uuid.Gen(), date: now},
			{uuid: uuid.Gen(), date: now.AddDate(0, 0, 7)},
		}

		userID := uuid.Gen()

		for _, input := range inputs {
			storage.CreateEvent(storage2.Event{ID: input.uuid, UserID: userID, StartAt: input.date, EndAt: input.date})
		}

		require.Len(t, storage.ListEventsWeek(userID, now.AddDate(0, 0, -7)), 1)
		require.Len(t, storage.ListEventsWeek(userID, now), 1)
		require.Len(t, storage.ListEventsWeek(userID, now.AddDate(0, 0, 7)), 1)
	})

	t.Run("ListEventsMonth", func(t *testing.T) {
		storage := storage2.NewMemStorage()
		now := time.Now()
		inputs := []struct {
			uuid uuid.UUID
			date time.Time
		}{
			{uuid: uuid.Gen(), date: now.AddDate(0, -1, 0)},
			{uuid: uuid.Gen(), date: now},
			{uuid: uuid.Gen(), date: now.AddDate(0, 1, 0)},
		}

		userID := uuid.Gen()

		for _, input := range inputs {
			storage.CreateEvent(storage2.Event{ID: input.uuid, UserID: userID, StartAt: input.date, EndAt: input.date})
		}

		require.Len(t, storage.ListEventsMonth(userID, now.AddDate(0, 1, 0)), 1)
		require.Len(t, storage.ListEventsMonth(userID, now), 1)
		require.Len(t, storage.ListEventsMonth(userID, now.AddDate(0, 1, 0)), 1)
	})
}
