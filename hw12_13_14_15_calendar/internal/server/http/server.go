package internalhttp

import (
	"context"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/server/http/actions"
	"net/http"
)

type Server struct {
	addr   string
	app    Application
	logger Logger

	server *http.Server
}

type Logger interface {
	Error(string)
	Warning(string)
	Info(string)
	Debug(string)
}

type Application interface {
}

func NewServer(logger Logger, app Application, host string, port string) *Server {
	return &Server{logger: logger, app: app, addr: host + ":" + port}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", actions.Ping)

	go func() {
		<-ctx.Done()

		err := s.Stop(ctx)
		if err != nil {
			return
		}
	}()

	s.server = &http.Server{Addr: s.addr, Handler: mux}

	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
