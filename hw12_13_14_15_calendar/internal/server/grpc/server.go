//go:generate protoc -I ../../../api/ --go_out=. --go-grpc_out=. ../../../api/EventService.proto
package internalgrpc

import (
	"context"
	"net"

	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/logger"
	"google.golang.org/grpc"
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

func (s *Server) CreateEventV1(ctx context.Context, event *EventV1) (*ResponseStatusV1, error) {
	if err := s.app.CreateEvent(
		ctx,
		event.GetId(),
		event.GetTitle(),
		event.GetDescription(),
		event.GetStartAt().AsTime(),
		event.GetEndAt().AsTime(),
		event.GetUserId(),
		event.GetRemindFor(),
	); err != nil {
		return &ResponseStatusV1{
			Code:    ResponseStatusV1_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &ResponseStatusV1{
		Code:    ResponseStatusV1_SUCCESS,
		Message: "",
	}, nil
}

func (s *Server) UpdateEventV1(ctx context.Context, event *EventV1) (*ResponseStatusV1, error) {
	if err := s.app.UpdateEvent(
		ctx,
		event.GetId(),
		event.GetTitle(),
		event.GetDescription(),
		event.GetStartAt().AsTime(),
		event.GetEndAt().AsTime(),
		event.GetUserId(),
		event.GetRemindFor(),
	); err != nil {
		return &ResponseStatusV1{
			Code:    ResponseStatusV1_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &ResponseStatusV1{
		Code:    ResponseStatusV1_SUCCESS,
		Message: "",
	}, nil
}

func (s *Server) DeleteEventV1(ctx context.Context, event *EventIdV1) (*ResponseStatusV1, error) {
	if err := s.app.DeleteEvent(ctx, event.GetId()); err != nil {
		return &ResponseStatusV1{
			Code:    ResponseStatusV1_FAILED,
			Message: err.Error(),
		}, nil
	}

	return &ResponseStatusV1{
		Code:    ResponseStatusV1_SUCCESS,
		Message: "",
	}, nil
}

func (s *Server) ListEventsDayV1(ctx context.Context, req *EventListRequestV1) (*EventResponseV1, error) {
	return s.listResponse(s.app.ListEventsDay(ctx, req.UserId, req.Date.AsTime())), nil
}

func (s *Server) ListEventsWeekV1(ctx context.Context, req *EventListRequestV1) (*EventResponseV1, error) {
	return s.listResponse(s.app.ListEventsWeek(ctx, req.UserId, req.Date.AsTime())), nil
}

func (s *Server) ListEventsMonthV1(ctx context.Context, req *EventListRequestV1) (*EventResponseV1, error) {
	return s.listResponse(s.app.ListEventsMonth(ctx, req.UserId, req.Date.AsTime())), nil
}

func (s *Server) listResponse(events []storage.Event) *EventResponseV1 {
	items := make([]*EventV1, len(events))
	for i, item := range events {
		item := item
		items[i] = &EventV1{
			Id:          item.ID.String(),
			Title:       item.Title,
			Description: &item.Description,
			StartAt:     timestamppb.New(item.StartAt),
			EndAt:       timestamppb.New(item.EndAt),
			UserId:      item.UserID.String(),
			RemindFor:   item.RemindFor,
		}
	}

	return &EventResponseV1{
		Items: items,
	}
}
