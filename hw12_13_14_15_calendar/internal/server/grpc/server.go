package internalgrpc

import (
	"context"
	"net"

	event "github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/api/stubs"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	addr   string
	app    app.Application
	logger *logger.Logger

	server *grpc.Server

	event.UnsafeCalendarServer
}

func NewServer(logger *logger.Logger, app app.Application, host string, port string) *Server {
	return &Server{logger: logger, app: app, addr: host + ":" + port}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		err := s.Stop(ctx)
		if err != nil {
			return
		}
	}()

	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.server = grpc.NewServer()
	event.RegisterCalendarServer(s.server, s)

	return s.server.Serve(listener)
}

func (s *Server) Stop(_ context.Context) error {
	s.server.GracefulStop()
	return nil
}

func (s *Server) CreateEventV1(ctx context.Context, v1 *event.EventV1) (*emptypb.Empty, error) {
	//s.app.CreateEvent(ctx, ...)
	panic("implement me")
}

func (s *Server) UpdateEventV1(ctx context.Context, v1 *event.EventV1) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) DeleteEventV1(ctx context.Context, v1 *event.EventIdV1) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) ListEventDayV1(ctx context.Context, timestamp *timestamppb.Timestamp) (*event.EventResponseV1, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) ListEventWeekV1(ctx context.Context, timestamp *timestamppb.Timestamp) (*event.EventResponseV1, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) ListEventMonthV1(ctx context.Context, timestamp *timestamppb.Timestamp) (*event.EventResponseV1, error) {
	//TODO implement me
	panic("implement me")
}
