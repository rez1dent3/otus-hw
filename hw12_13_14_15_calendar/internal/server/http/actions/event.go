package actions

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/server/http/requests"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/server/http/responses"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

type EventEnt struct {
	app    app.Application
	logger *logger.Logger
}

func NewEventEnt(application app.Application, logger *logger.Logger) *EventEnt {
	return &EventEnt{app: application, logger: logger}
}

// createEvent POST /events {...}.
func (e *EventEnt) createEvent(writer http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		e.error(writer, err)
		return
	}

	eventReq := requests.EventRequest{}
	if err = json.Unmarshal(body, &eventReq); err != nil {
		e.error(writer, err)
		return
	}

	if err = e.app.CreateEvent(
		req.Context(),
		eventReq.ID,
		eventReq.Title,
		eventReq.Description,
		eventReq.StartAt,
		eventReq.EndAt,
		eventReq.UserID,
		eventReq.RemindFor,
	); err != nil {
		e.error(writer, err)
		return
	}

	e.json(writer, struct{}{})
}

// updateEvent PUT /events {...}.
func (e *EventEnt) updateEvent(writer http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		e.error(writer, err)
		return
	}

	eventReq := requests.EventRequest{}
	if err = json.Unmarshal(body, &eventReq); err != nil {
		e.error(writer, err)
		return
	}

	if err = e.app.UpdateEvent(
		req.Context(),
		eventReq.ID,
		eventReq.Title,
		eventReq.Description,
		eventReq.StartAt,
		eventReq.EndAt,
		eventReq.UserID,
		eventReq.RemindFor,
	); err != nil {
		e.error(writer, err)
		return
	}

	e.json(writer, struct{}{})
}

// deleteEvent DELETE /events?id=:id.
func (e *EventEnt) deleteEvent(writer http.ResponseWriter, req *http.Request) {
	if err := e.app.DeleteEvent(req.Context(), req.URL.Query().Get("id")); err != nil {
		e.error(writer, err)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

// listEvent GET /events?type=:type&userId=:userId&timestamp=:timestamp.
func (e *EventEnt) listEvent(writer http.ResponseWriter, req *http.Request) {
	timestampString := req.URL.Query().Get("timestamp")
	timestamp, err := strconv.ParseInt(timestampString, 10, 64)
	if err != nil {
		e.error(writer, err)
		return
	}

	listEvent := requests.EventListRequest{
		UserID: req.URL.Query().Get("userId"),
		Date:   time.Unix(timestamp, 0),
	}

	switch req.URL.Query().Get("type") {
	case "ListEventsDay":
		e.json(writer, e.list(e.app.ListEventsDay(req.Context(), listEvent.UserID, listEvent.Date)))
	case "ListEventsWeek":
		e.json(writer, e.list(e.app.ListEventsWeek(req.Context(), listEvent.UserID, listEvent.Date)))
	case "ListEventsMonth":
		e.json(writer, e.list(e.app.ListEventsMonth(req.Context(), listEvent.UserID, listEvent.Date)))
	default:
		http.NotFound(writer, req)
	}
}

func (e *EventEnt) error(w http.ResponseWriter, err error) {
	e.logger.Error(err.Error())

	resp := responses.ErrorResponse{Message: err.Error()}
	w.WriteHeader(http.StatusBadRequest)
	e.json(w, resp)
}

func (e *EventEnt) list(events []storage.Event) []responses.EventResponse {
	result := make([]responses.EventResponse, len(events))
	i := 0
	for _, item := range events {
		item := item
		result[i] = responses.EventResponse{
			ID:          item.ID.String(),
			Title:       item.Title,
			Description: &item.Description,
			StartAt:     item.StartAt,
			EndAt:       item.EndAt,
			UserID:      item.UserID.String(),
			RemindFor:   item.RemindFor,
		}
		i++
	}
	return result
}

func (e *EventEnt) json(w http.ResponseWriter, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		e.logger.Error(err.Error())
	}
}

func (e *EventEnt) HandleFunc(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		e.listEvent(writer, req)
		return
	case http.MethodPost:
		e.createEvent(writer, req)
		return
	case http.MethodPut:
		e.updateEvent(writer, req)
		return
	case http.MethodDelete:
		e.deleteEvent(writer, req)
		return
	default:
		http.NotFound(writer, req)
	}
}
