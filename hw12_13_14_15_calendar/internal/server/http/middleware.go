package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

type MiddlewareInterface interface {
	Handle(http.Handler) http.Handler
}

type responseWriteStatusDecorator struct {
	http.ResponseWriter
	status int
}

type LoggerMiddleware struct {
	logger Logger
}

func (w *responseWriteStatusDecorator) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriteStatusDecorator) GetStatus() int {
	return w.status
}

func NewLoggerMiddleware(logger Logger) MiddlewareInterface {
	return &LoggerMiddleware{logger: logger}
}

func (m *LoggerMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w2 := responseWriteStatusDecorator{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		now := time.Now()
		next.ServeHTTP(&w2, r)
		latency := time.Since(now)

		m.logger.Info(fmt.Sprintf(
			"%s [%s] %s %s %s %d %d %s",
			r.Header.Get("X-Forwarded-For"),
			now.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.RequestURI,
			r.Proto,
			w2.GetStatus(),
			latency.Microseconds(),
			r.Header.Get("User-Agent"),
		))
	})
}
