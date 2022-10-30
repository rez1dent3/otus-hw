//go:generate protoc -I ../../../api/ --go_out=. --go-grpc_out=. ../../../api/EventService.proto
package internalgrpc

import (
	"context"
	"net"

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

	UnsafeCalendarServer
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

	s.server = grpc.NewServer(grpc.UnaryInterceptor(s.RequestInterceptor))
	RegisterCalendarServer(s.server, s)

	return s.server.Serve(listener)
}

func (s *Server) Stop(_ context.Context) error {
	s.server.GracefulStop()
	return nil
}

func (s *Server) CreateEventV1(context.Context, *EventV1) (*emptypb.Empty, error) {
	//s.app.CreateEvent(ctx, ...)
	panic("implement me")
}

func (s *Server) UpdateEventV1(context.Context, *EventV1) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) DeleteEventV1(context.Context, *EventIdV1) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) ListEventsDayV1(context.Context, *timestamppb.Timestamp) (*EventResponseV1, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) ListEventsWeekV1(context.Context, *timestamppb.Timestamp) (*EventResponseV1, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) ListEventsMonthV1(context.Context, *timestamppb.Timestamp) (*EventResponseV1, error) {
	//TODO implement me
	panic("implement me")
}
