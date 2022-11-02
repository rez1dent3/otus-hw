package internalgrpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/app"
	internalgrpc "github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type nullWriter struct{}

func (w nullWriter) Write(_ []byte) (n int, err error) {
	return 0, nil
}

func createServer() *internalgrpc.Server {
	mem := storage.NewMemStorage()
	log := logger.New("off", &nullWriter{})
	application := app.New(log, mem)

	return internalgrpc.NewServer(log, application, "", "")
}

func newEvent() *internalgrpc.EventV1 {
	eventID := uuid.Gen()
	userID := uuid.Gen()

	return &internalgrpc.EventV1{
		Id:          eventID.String(),
		Title:       "test",
		Description: nil,
		StartAt:     timestamppb.New(time.Unix(1667323499, 10)),
		EndAt:       timestamppb.New(time.Unix(1667399693, 10)),
		UserId:      userID.String(),
		RemindFor:   nil,
	}
}

func TestServer_CreateEventV1(t *testing.T) {
	t.Parallel()

	t.Run("CreateEventV1", func(t *testing.T) {
		server := createServer()
		event := newEvent()
		_, err := server.CreateEventV1(context.Background(), event)
		require.NoError(t, err)

		resp, err := server.ListEventsDayV1(context.Background(), &internalgrpc.EventListRequestV1{
			UserId: event.UserId,
			Date:   timestamppb.New(time.Unix(1667323499, 10)),
		})

		require.NoError(t, err)
		require.Len(t, resp.Items, 1)
		require.Equal(t, event.Id, resp.Items[0].Id)
		require.Equal(t, event.UserId, resp.Items[0].UserId)
	})

	t.Run("CreateEventV1.ErrDuplicate", func(t *testing.T) {
		server := createServer()
		event := newEvent()

		_, err := server.CreateEventV1(context.Background(), event)
		require.NoError(t, err)

		_, err = server.CreateEventV1(context.Background(), event)
		require.ErrorIs(t, app.ErrCreateEvent, err)
	})
}

func TestServer_DeleteEventV1(t *testing.T) {
	t.Parallel()

	t.Run("DeleteEventV1", func(t *testing.T) {
		server := createServer()
		event := newEvent()

		_, err := server.CreateEventV1(context.Background(), event)
		require.NoError(t, err)

		_, err = server.DeleteEventV1(context.Background(), &internalgrpc.EventIdV1{
			Id: event.Id,
		})
		require.NoError(t, err)

		_, err = server.CreateEventV1(context.Background(), event)
		require.NoError(t, err)
	})
}

func TestServer_UpdateEventV1(t *testing.T) {
	t.Parallel()

	t.Run("UpdateEventV1", func(t *testing.T) {
		server := createServer()
		event := newEvent()

		_, err := server.CreateEventV1(context.Background(), event)
		require.NoError(t, err)

		userID := uuid.Gen()

		{
			event := event
			event.UserId = userID.String()
			_, err = server.UpdateEventV1(context.Background(), event)
			require.NoError(t, err)
		}

		resp, err := server.ListEventsDayV1(context.Background(), &internalgrpc.EventListRequestV1{
			UserId: userID.String(),
			Date:   timestamppb.New(time.Unix(1667323499, 10)),
		})

		require.NoError(t, err)
		require.Len(t, resp.Items, 1)
		require.Equal(t, event.Id, resp.Items[0].Id)
		require.Equal(t, userID.String(), resp.Items[0].UserId)
	})
}

func TestServer_ListEventsBy(t *testing.T) {
	t.Parallel()

	t.Run("ListEventsDayV1", func(t *testing.T) {
		server := createServer()
		event := newEvent()

		_, err := server.CreateEventV1(context.Background(), event)
		require.NoError(t, err)

		resp, err := server.ListEventsDayV1(context.Background(), &internalgrpc.EventListRequestV1{
			UserId: event.UserId,
			Date:   timestamppb.New(time.Unix(1667433600, 0)),
		})

		require.NoError(t, err)
		require.Len(t, resp.Items, 0)

		resp, err = server.ListEventsDayV1(context.Background(), &internalgrpc.EventListRequestV1{
			UserId: event.UserId,
			Date:   timestamppb.New(time.Unix(1667323499, 0)),
		})

		require.NoError(t, err)
		require.Len(t, resp.Items, 1)
	})

	t.Run("ListEventsWeekV1", func(t *testing.T) {
		server := createServer()
		event := newEvent()

		_, err := server.CreateEventV1(context.Background(), event)
		require.NoError(t, err)

		resp, err := server.ListEventsWeekV1(context.Background(), &internalgrpc.EventListRequestV1{
			UserId: event.UserId,
			Date:   timestamppb.New(time.Unix(1667174399, 0)),
		})

		require.NoError(t, err)
		require.Len(t, resp.Items, 0)

		resp, err = server.ListEventsWeekV1(context.Background(), &internalgrpc.EventListRequestV1{
			UserId: event.UserId,
			Date:   timestamppb.New(time.Unix(1667779200, 0)),
		})

		require.NoError(t, err)
		require.Len(t, resp.Items, 0)

		resp, err = server.ListEventsWeekV1(context.Background(), &internalgrpc.EventListRequestV1{
			UserId: event.UserId,
			Date:   timestamppb.New(time.Unix(1667174400, 0)),
		})

		require.NoError(t, err)
		require.Len(t, resp.Items, 1)
	})

	t.Run("ListEventsMonthV1", func(t *testing.T) {
		server := createServer()
		event := newEvent()

		_, err := server.CreateEventV1(context.Background(), event)
		require.NoError(t, err)

		resp, err := server.ListEventsMonthV1(context.Background(), &internalgrpc.EventListRequestV1{
			UserId: event.UserId,
			Date:   timestamppb.New(time.Unix(1667260799, 0)),
		})

		require.NoError(t, err)
		require.Len(t, resp.Items, 0)

		resp, err = server.ListEventsMonthV1(context.Background(), &internalgrpc.EventListRequestV1{
			UserId: event.UserId,
			Date:   timestamppb.New(time.Unix(1669852800, 0)),
		})

		require.NoError(t, err)
		require.Len(t, resp.Items, 0)

		resp, err = server.ListEventsMonthV1(context.Background(), &internalgrpc.EventListRequestV1{
			UserId: event.UserId,
			Date:   timestamppb.New(time.Unix(1667260800, 0)),
		})

		require.NoError(t, err)
		require.Len(t, resp.Items, 1)
	})
}
