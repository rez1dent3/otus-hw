package internalhttp_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/app"
	internalhttp "github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/server/http/requests"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
	"github.com/stretchr/testify/require"
)

type nullWriter struct{}

func (w nullWriter) Write(_ []byte) (n int, err error) {
	return 0, nil
}

func createServer() *internalhttp.Server {
	mem := storage.NewMemStorage()
	log := logger.New("off", &nullWriter{})
	application := app.New(log, mem)

	return internalhttp.NewServer(log, application, "", "")
}

func TestServer_Infra(t *testing.T) {
	t.Parallel()

	t.Run("Healthcheck", func(t *testing.T) {
		server := httptest.NewServer(createServer().HTTPHandler())
		defer server.Close()

		res, err := http.Get(fmt.Sprintf("%s/health", server.URL))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, http.StatusOK)

		defer func() {
			_ = res.Body.Close()
		}()

		body, err := io.ReadAll(res.Body)
		require.Equal(t, "ok", string(body))
	})
}

func TestServer_CreateEvent(t *testing.T) {
	t.Parallel()

	t.Run("CreateEvent", func(t *testing.T) {
		server := httptest.NewServer(createServer().HTTPHandler())
		defer server.Close()

		eventID := uuid.Gen()
		userID := uuid.Gen()

		createEvent := requests.EventRequest{
			ID:      eventID.String(),
			StartAt: time.Now(),
			EndAt:   time.Now(),
			UserID:  userID.String(),
		}

		body, err := json.Marshal(createEvent)
		require.NoError(t, err)

		bodyReader := strings.NewReader(string(body))

		res, err := http.Post(fmt.Sprintf("%s/events", server.URL), "application/json", bodyReader)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

		defer func() {
			_ = res.Body.Close()
		}()

		result, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, "{}\n", string(result))
	})

	t.Run("CreateEvent.Duplicate", func(t *testing.T) {
		server := httptest.NewServer(createServer().HTTPHandler())
		defer server.Close()

		eventID := uuid.Gen()
		userID := uuid.Gen()

		createEvent := requests.EventRequest{
			ID:      eventID.String(),
			StartAt: time.Now(),
			EndAt:   time.Now(),
			UserID:  userID.String(),
		}

		body, err := json.Marshal(createEvent)
		require.NoError(t, err)

		bodyReader := strings.NewReader(string(body))

		res, err := http.Post(fmt.Sprintf("%s/events", server.URL), "application/json", bodyReader)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

		defer func() {
			_ = res.Body.Close()
		}()

		result, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, "{}\n", string(result))

		bodyReader = strings.NewReader(string(body))
		res, err = http.Post(fmt.Sprintf("%s/events", server.URL), "application/json", bodyReader)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)

		defer func() {
			_ = res.Body.Close()
		}()

		body, err = io.ReadAll(res.Body)
		require.Equal(t, "{\"message\":\"create event error\"}\n", string(body))
	})

	t.Run("DeleteEvent", func(t *testing.T) {
		server := httptest.NewServer(createServer().HTTPHandler())
		defer server.Close()

		eventID := uuid.Gen()
		userID := uuid.Gen()

		createEvent := requests.EventRequest{
			ID:      eventID.String(),
			StartAt: time.Now(),
			EndAt:   time.Now(),
			UserID:  userID.String(),
		}

		body, err := json.Marshal(createEvent)
		require.NoError(t, err)

		bodyReader := strings.NewReader(string(body))

		res, err := http.Post(fmt.Sprintf("%s/events", server.URL), "application/json", bodyReader)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

		defer func() {
			_ = res.Body.Close()
		}()

		result, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, "{}\n", string(result))

		client := http.Client{Timeout: time.Second}
		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("%s/events?id="+eventID.String(), server.URL), nil)
		require.NoError(t, err)

		res, err = client.Do(req)
		require.NoError(t, err)

		defer func() {
			_ = res.Body.Close()
		}()

		body, err = io.ReadAll(res.Body)
		require.Equal(t, "", string(body))
		require.Equal(t, http.StatusNoContent, res.StatusCode)
	})
}
