package memorystorage

import (
	storage2 "github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestStorage_CRUD(t *testing.T) {
	t.Run("CreateEvent", func(t *testing.T) {
		storage := New()
		event := storage2.Event{}

		require.Len(t, storage.ListEventsMonth(event.StartAt), 0)
		success := storage.CreateEvent(storage2.Event{})

		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(event.StartAt), 1)
	})

	t.Run("UpdateEvent", func(t *testing.T) {
		storage := New()
		event1 := storage2.Event{}
		loss := storage.UpdateEvent(event1.ID, event1)
		require.False(t, loss)

		success := storage.CreateEvent(event1)
		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(event1.StartAt), 1)

		event2 := storage2.Event{ID: event1.ID, Title: "hello"}
		success = storage.UpdateEvent(event2.ID, event2)
		require.Len(t, storage.ListEventsMonth(event2.StartAt), 1)
		require.True(t, success)
	})

	t.Run("DeleteEvent.notExists", func(t *testing.T) {
		storage := New()
		event1 := storage2.Event{}
		require.Len(t, storage.ListEventsMonth(event1.StartAt), 0)
		success := storage.DeleteEvent(event1.ID)
		require.True(t, success)
	})

	t.Run("DeleteEvent.exists", func(t *testing.T) {
		storage := New()
		event1 := storage2.Event{}
		success := storage.CreateEvent(event1)
		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(event1.StartAt), 1)

		success = storage.DeleteEvent(event1.ID)
		require.True(t, success)
		require.Len(t, storage.ListEventsMonth(event1.StartAt), 0)
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

		storage := New()
		for _, input := range inputs {
			storage.CreateEvent(storage2.Event{ID: uuid.FromString(input.uuid), StartAt: input.date, EndAt: input.date})
		}

		require.Len(t, storage.ListEventsDay(now.AddDate(0, 0, -1)), 1)
		require.Len(t, storage.ListEventsDay(now), 1)
		require.Len(t, storage.ListEventsDay(now.AddDate(0, 0, 1)), 1)
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

		storage := New()
		for _, input := range inputs {
			storage.CreateEvent(storage2.Event{ID: uuid.FromString(input.uuid), StartAt: input.date, EndAt: input.date})
		}

		require.Len(t, storage.ListEventsWeek(now.AddDate(0, 0, -7)), 1)
		require.Len(t, storage.ListEventsWeek(now), 1)
		require.Len(t, storage.ListEventsWeek(now.AddDate(0, 0, 7)), 1)
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

		storage := New()
		for _, input := range inputs {
			storage.CreateEvent(storage2.Event{ID: uuid.FromString(input.uuid), StartAt: input.date, EndAt: input.date})
		}

		require.Len(t, storage.ListEventsMonth(now.AddDate(0, 1, 0)), 1)
		require.Len(t, storage.ListEventsMonth(now), 1)
		require.Len(t, storage.ListEventsMonth(now.AddDate(0, 1, 0)), 1)
	})

}
