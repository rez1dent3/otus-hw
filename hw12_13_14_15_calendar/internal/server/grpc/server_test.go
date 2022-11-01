package internalgrpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/app"
	internalgrpc "github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type nullWriter struct{}

func (w nullWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

func createServer() *internalgrpc.Server {
	mem := storage.NewMemStorage()
	log := logger.New("off", &nullWriter{})
	application := app.New(log, mem)

	return internalgrpc.NewServer(log, application, "", "")
}

func TestServer(t *testing.T) {
	event := &internalgrpc.EventV1{
		Id:          "7feb6c19-f557-4536-ab2c-a8798e628909",
		Title:       "test",
		Description: nil,
		StartAt:     timestamppb.New(time.Unix(1667323499, 10)),
		EndAt:       timestamppb.New(time.Unix(1667399693, 10)),
		UserId:      "388d1e70-38e8-45be-a889-cc02671b1409",
		RemindFor:   nil,
	}

	t.Run("create duplicate", func(t *testing.T) {
		s := createServer()
		_, err := s.CreateEventV1(context.Background(), event)
		require.NoError(t, err)

		_, err = s.CreateEventV1(context.Background(), event)
		require.ErrorIs(t, app.ErrCreateEvent, err)
	})
}
