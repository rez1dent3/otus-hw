package storage_test

import (
	"context"
	"testing"
	"time"

	_ "github.com/lib/pq"
	storage2 "github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
	"github.com/stretchr/testify/require"
)

const dsn = "host=127.0.0.1 port=5432 user=calendar_user password=calendar_pass dbname=calendar sslmode=disable"

func TestPgStorage_Create(t *testing.T) {
	t.Skip()

	storage := storage2.NewPgStorage(dsn)
	storage.Connect(context.Background())
	defer storage.Close()

	t.Run("CreateEvent.notExists", func(t *testing.T) {
		event := storage2.Event{ID: uuid.Gen(), UserID: uuid.Gen()}

		require.Len(t, storage.ListEventsMonth(context.Background(), event.UserID, event.StartAt), 0)
		success, err := storage.CreateEvent(context.Background(), event)
		require.NoError(t, err)

		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(context.Background(), event.UserID, event.StartAt), 1)
	})

	t.Run("CreateEvent.exists", func(t *testing.T) {
		event := storage2.Event{ID: uuid.Gen(), UserID: uuid.Gen()}

		require.Len(t, storage.ListEventsMonth(context.Background(), event.UserID, event.StartAt), 0)
		success, err := storage.CreateEvent(context.Background(), event)
		require.NoError(t, err)

		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(context.Background(), event.UserID, event.StartAt), 1)

		loss, err := storage.CreateEvent(context.Background(), event)
		require.ErrorIs(t, err, storage2.ErrUnableDuplicate)
		require.False(t, loss)
		require.Len(t, storage.ListEventsMonth(context.Background(), event.UserID, event.StartAt), 1)
	})
}

func TestPgStorage_Update(t *testing.T) {
	t.Skip()

	storage := storage2.NewPgStorage(dsn)
	storage.Connect(context.Background())
	defer storage.Close()

	t.Run("UpdateEvent", func(t *testing.T) {
		event1 := storage2.Event{ID: uuid.Gen(), UserID: uuid.Gen()}
		loss, err := storage.UpdateEvent(context.Background(), event1.ID, event1)
		require.ErrorIs(t, err, storage2.ErrNotFound)
		require.False(t, loss)

		success, err := storage.CreateEvent(context.Background(), event1)
		require.NoError(t, err)
		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(context.Background(), event1.UserID, event1.StartAt), 1)

		event2 := storage2.Event{ID: event1.ID, UserID: uuid.Gen(), Title: "hello"}
		success, err = storage.UpdateEvent(context.Background(), event2.ID, event2)
		require.NoError(t, err)
		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(context.Background(), event2.UserID, event2.StartAt), 1)
	})
}

func TestPgStorage_Delete(t *testing.T) {
	t.Skip()

	storage := storage2.NewPgStorage(dsn)
	storage.Connect(context.Background())
	defer storage.Close()

	t.Run("DeleteEvent.notExists", func(t *testing.T) {
		event1 := storage2.Event{ID: uuid.Gen(), UserID: uuid.Gen()}
		require.Len(t, storage.ListEventsMonth(context.Background(), event1.UserID, event1.StartAt), 0)
		success, err := storage.DeleteEvent(context.Background(), event1.ID)
		require.ErrorIs(t, err, storage2.ErrNotFound)
		require.False(t, success)
	})

	t.Run("DeleteEvent.exists", func(t *testing.T) {
		event1 := storage2.Event{ID: uuid.Gen(), UserID: uuid.Gen()}
		success, err := storage.CreateEvent(context.Background(), event1)
		require.NoError(t, err)
		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(context.Background(), event1.UserID, event1.StartAt), 1)

		success, err = storage.DeleteEvent(context.Background(), event1.ID)
		require.NoError(t, err)
		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(context.Background(), event1.UserID, event1.StartAt), 0)
	})
}

func TestPgStorage_List(t *testing.T) {
	t.Skip()

	storage := storage2.NewPgStorage(dsn)
	storage.Connect(context.Background())
	defer storage.Close()

	t.Run("list.filterByUserID", func(t *testing.T) {
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
				context.Background(),
				storage2.Event{
					ID:      input.eventID,
					UserID:  input.userID,
					StartAt: input.date,
					EndAt:   input.date,
				})
		}

		require.Len(t, storage.ListEventsDay(context.Background(), user1, now.AddDate(0, 0, -1)), 1)
		require.Len(t, storage.ListEventsDay(context.Background(), user1, now), 1)
		require.Len(t, storage.ListEventsDay(context.Background(), user1, now.AddDate(0, 0, 1)), 0)

		require.Len(t, storage.ListEventsDay(context.Background(), user2, now.AddDate(0, 0, -1)), 0)
		require.Len(t, storage.ListEventsDay(context.Background(), user2, now), 0)
		require.Len(t, storage.ListEventsDay(context.Background(), user2, now.AddDate(0, 0, 1)), 1)
	})

	t.Run("ListEventsDay", func(t *testing.T) {
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
			storage.CreateEvent(
				context.Background(),
				storage2.Event{ID: input.uuid, UserID: userID, StartAt: input.date, EndAt: input.date})
		}

		require.Len(t, storage.ListEventsDay(context.Background(), userID, now.AddDate(0, 0, -1)), 1)
		require.Len(t, storage.ListEventsDay(context.Background(), userID, now), 1)
		require.Len(t, storage.ListEventsDay(context.Background(), userID, now.AddDate(0, 0, 1)), 1)
	})

	t.Run("ListEventsWeek", func(t *testing.T) {
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
			storage.CreateEvent(
				context.Background(),
				storage2.Event{ID: input.uuid, UserID: userID, StartAt: input.date, EndAt: input.date})
		}

		require.Len(t, storage.ListEventsWeek(context.Background(), userID, now.AddDate(0, 0, -7)), 1)
		require.Len(t, storage.ListEventsWeek(context.Background(), userID, now), 1)
		require.Len(t, storage.ListEventsWeek(context.Background(), userID, now.AddDate(0, 0, 7)), 1)
	})

	t.Run("ListEventsMonth", func(t *testing.T) {
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
			storage.CreateEvent(
				context.Background(),
				storage2.Event{ID: input.uuid, UserID: userID, StartAt: input.date, EndAt: input.date})
		}

		require.Len(t, storage.ListEventsMonth(context.Background(), userID, now.AddDate(0, 1, 0)), 1)
		require.Len(t, storage.ListEventsMonth(context.Background(), userID, now), 1)
		require.Len(t, storage.ListEventsMonth(context.Background(), userID, now.AddDate(0, 1, 0)), 1)
	})
}

func TestPgStorage_Notify(t *testing.T) {
	t.Skip()

	storage := storage2.NewPgStorage(dsn)
	storage.Connect(context.Background())
	defer storage.Close()

	t.Run("ListToSendNotifies", func(t *testing.T) {
		now := time.Now().UTC()
		inputs := []struct {
			uuid uuid.UUID
			date time.Time
			dur  time.Duration
		}{
			{uuid: uuid.Gen(), date: now, dur: time.Hour},
		}

		for _, input := range inputs {
			_, err := storage.CreateEvent(
				context.Background(),
				storage2.Event{
					ID:        input.uuid,
					UserID:    uuid.Gen(),
					StartAt:   input.date,
					EndAt:     input.date.Add(24 * time.Hour),
					RemindFor: (*storage2.Duration)(&input.dur)})

			require.NoError(t, err)
		}

		notifies, err := storage.ListToSendNotifies(context.Background(), now.AddDate(0, 0, -1))
		require.NoError(t, err)
		require.Len(t, notifies, 0)

		notifies, err = storage.ListToSendNotifies(context.Background(), now.Add(-(2 * time.Hour)))
		require.NoError(t, err)
		require.Len(t, notifies, 0)

		notifies, err = storage.ListToSendNotifies(context.Background(), now.Add(2*time.Hour))
		require.NoError(t, err)
		require.Len(t, notifies, 1)

		list := storage.ListEventsDay(context.Background(), notifies[0].UserID, notifies[0].StartAt)
		event := list[notifies[0].EventID]

		require.NotNil(t, event.RemindFor)
		require.Equal(t, time.Hour, (time.Duration)(*event.RemindFor))

		notifies, err = storage.ListToSendNotifies(context.Background(), now)
		require.NoError(t, err)
		require.Len(t, notifies, 1)

		for _, input := range inputs {
			require.NoError(t, storage.MarkInQueue(context.Background(), input.uuid))
		}

		notifies, err = storage.ListToSendNotifies(context.Background(), now)
		require.NoError(t, err)
		require.Len(t, notifies, 0)
	})
}
