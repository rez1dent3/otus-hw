package internalhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/server/http/actions"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

type Server struct {
	addr   string
	app    app.Application
	logger *logger.Logger

	server *http.Server
}

func NewServer(logger *logger.Logger, app app.Application, host string, port string) *Server {
	return &Server{logger: logger, app: app, addr: host + ":" + port}
}

func (s *Server) HTTPHandler() http.Handler {
	stack := actions.NewEventEnt(s.app, s.logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", actions.Ping)
	mux.HandleFunc("/events", stack.HandleFunc)

	return mux
}

func (s *Server) Start(ctx context.Context) error {
	mux := s.HTTPHandler()

	loggerMiddleware := NewLoggerMiddleware(s.logger)
	s.server = &http.Server{Addr: s.addr, Handler: loggerMiddleware.Handle(mux), ReadHeaderTimeout: time.Second}

	go func() {
		<-ctx.Done()

		err := s.Stop(ctx)
		if err != nil {
			return
		}
	}()

	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
