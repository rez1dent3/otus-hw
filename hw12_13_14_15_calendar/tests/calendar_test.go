package tests__test

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/server/http/requests"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/server/http/responses"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
	"github.com/stretchr/testify/suite"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	HTTP serverConf
	Grpc serverConf
}

type serverConf struct {
	Host string
	Port string
}

func NewConfig(reader io.Reader) (*Config, error) {
	cfg := Config{}
	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

type httpTestSuite struct {
	suite.Suite

	cfg *Config

	event *requests.EventRequest

	startAt   time.Time
	endAt     time.Time
	remindFor *time.Duration

	eventID uuid.UUID
	userID  uuid.UUID
}

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/config.yaml", "Path to configuration file")
}

func TestHttpTestSuite(t *testing.T) {
	suite.Run(t, &httpTestSuite{})
}

func (s *httpTestSuite) SetupSuite() {
	flag.Parse()

	file, err := os.Open(configFile)
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		_ = file.Close()
	}()

	cfg, err := NewConfig(file)
	if err != nil {
		log.Println(err)
		return
	}

	s.cfg = cfg
}

func (s *httpTestSuite) SetupTest() {
	s.eventID = uuid.Gen()
	s.userID = uuid.Gen()
	s.startAt = time.Now().Add(4 * time.Hour).UTC()
	s.endAt = time.Now().Add(24 * time.Hour).UTC()

	remindFor := 4 * time.Hour
	s.remindFor = &remindFor

	s.event = &requests.EventRequest{
		ID:    s.eventID.String(),
		Title: "Lorem Ipsum",
		Description: `Lorem Ipsum is simply dummy text of the printing and typesetting industry. 
				Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, 
				when an unknown printer took a galley of type and scrambled it to make a type specimen book. 
				It has survived not only five centuries, but also the leap into electronic typesetting, 
				remaining essentially unchanged. 

				It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, 
				and more recently with desktop publishing software like Aldus PageMaker 
				including versions of Lorem Ipsum.`,

		StartAt:   s.startAt,
		EndAt:     s.endAt,
		RemindFor: s.remindFor,

		UserID: s.userID.String(),
	}
}

func (s *httpTestSuite) TearDownTest() {
	s.remindFor = nil
	s.event = nil
}

func (s *httpTestSuite) TestCheckCreateEvent() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reqBody, err := json.Marshal(s.event)
	s.NoError(err)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("http://%s:%s/events", s.cfg.HTTP.Host, s.cfg.HTTP.Port),
		bytes.NewReader(reqBody))

	s.NoError(err)

	client := http.Client{}

	resp, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.NoError(resp.Body.Close())

	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"http://%s:%s/events?type=ListEventsDay&userId=%s&timestamp=%d",
			s.cfg.HTTP.Host,
			s.cfg.HTTP.Port,
			s.userID.String(),
			s.startAt.Unix(),
		),
		nil)

	resp, err = client.Do(req)
	s.NoError(err)

	respBody, err := io.ReadAll(resp.Body)
	s.NoError(resp.Body.Close())
	s.Equal(http.StatusOK, resp.StatusCode)

	var events []responses.EventResponse
	err = json.Unmarshal(respBody, &events)
	s.NoError(err)

	s.Equal(s.event.ID, events[0].ID)
	s.Equal(s.event.Title, events[0].Title)
	s.Equal(s.event.Description, *events[0].Description)
	s.Equal(s.event.StartAt.Truncate(time.Second), events[0].StartAt.Truncate(time.Second))
	s.Equal(s.event.EndAt.Truncate(time.Second), events[0].EndAt.Truncate(time.Second))
	s.Equal(s.event.UserID, events[0].UserID)
	s.Equal(*s.event.RemindFor, *events[0].RemindFor)
	s.Equal(false, events[0].IsDispatched)
}

func (s *httpTestSuite) TestCheckDuplicateCreateEvent() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reqBody, err := json.Marshal(s.event)
	s.NoError(err)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("http://%s:%s/events", s.cfg.HTTP.Host, s.cfg.HTTP.Port),
		bytes.NewReader(reqBody))

	s.NoError(err)

	client := http.Client{}

	resp, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.NoError(resp.Body.Close())

	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("http://%s:%s/events", s.cfg.HTTP.Host, s.cfg.HTTP.Port),
		bytes.NewReader(reqBody))

	s.NoError(err)

	resp, err = client.Do(req)
	s.NoError(err)

	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	s.NoError(err)

	s.Equal("{\"message\":\"unable to duplicate\"}\n", string(respBody))
	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *httpTestSuite) TestCheckUpdateEvent() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reqBody, err := json.Marshal(s.event)
	s.NoError(err)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("http://%s:%s/events", s.cfg.HTTP.Host, s.cfg.HTTP.Port),
		bytes.NewReader(reqBody))

	s.NoError(err)

	client := http.Client{}

	resp, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.NoError(resp.Body.Close())

	newUserID := uuid.Gen()
	newEvent := requests.EventRequest{
		ID:          s.event.ID,
		Title:       "Title",
		Description: "Description",
		StartAt:     s.event.EndAt,
		EndAt:       s.endAt.Add(*s.remindFor),
		UserID:      newUserID.String(),
		RemindFor:   nil,
	}

	reqBody, err = json.Marshal(newEvent)
	s.NoError(err)

	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf("http://%s:%s/events", s.cfg.HTTP.Host, s.cfg.HTTP.Port),
		bytes.NewReader(reqBody))

	s.NoError(err)

	resp, err = client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.NoError(resp.Body.Close())

	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"http://%s:%s/events?type=ListEventsDay&userId=%s&timestamp=%d",
			s.cfg.HTTP.Host,
			s.cfg.HTTP.Port,
			newUserID.String(),
			s.endAt.Unix(),
		),
		nil)

	resp, err = client.Do(req)
	s.NoError(err)

	respBody, err := io.ReadAll(resp.Body)
	s.NoError(resp.Body.Close())
	s.Equal(http.StatusOK, resp.StatusCode)

	var events []responses.EventResponse
	err = json.Unmarshal(respBody, &events)
	s.NoError(err)

	s.Equal(s.eventID.String(), events[0].ID)
	s.Equal("Title", events[0].Title)
	s.Equal("Description", *events[0].Description)
	s.Equal(s.event.EndAt.Truncate(time.Second), events[0].StartAt.Truncate(time.Second))
	s.Equal(s.event.EndAt.Add(*s.event.RemindFor).Truncate(time.Second), events[0].EndAt.Truncate(time.Second))
	s.Equal(newUserID.String(), events[0].UserID)
	s.Equal(false, events[0].IsDispatched)
}

func (s *httpTestSuite) TestCheckUpdateEventNotFound() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reqBody, err := json.Marshal(s.event)
	s.NoError(err)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf("http://%s:%s/events", s.cfg.HTTP.Host, s.cfg.HTTP.Port),
		bytes.NewReader(reqBody))

	s.NoError(err)

	client := http.Client{}

	resp, err := client.Do(req)
	s.NoError(err)

	defer func() {
		_ = resp.Body.Close()
	}()

	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *httpTestSuite) TestCheckDeleteEvent() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reqBody, err := json.Marshal(s.event)
	s.NoError(err)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("http://%s:%s/events", s.cfg.HTTP.Host, s.cfg.HTTP.Port),
		bytes.NewReader(reqBody))

	s.NoError(err)

	client := http.Client{}

	resp, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.NoError(resp.Body.Close())

	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf(
			"http://%s:%s/events?id=%s",
			s.cfg.HTTP.Host,
			s.cfg.HTTP.Port,
			s.eventID.String(),
		),
		nil)

	resp, err = client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusNoContent, resp.StatusCode)
	s.NoError(resp.Body.Close())

	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"http://%s:%s/events?type=ListEventsDay&userId=%s&timestamp=%d",
			s.cfg.HTTP.Host,
			s.cfg.HTTP.Port,
			s.userID.String(),
			s.startAt.Unix(),
		),
		nil)

	resp, err = client.Do(req)
	s.NoError(err)

	respBody, err := io.ReadAll(resp.Body)
	s.NoError(resp.Body.Close())
	s.Equal(http.StatusOK, resp.StatusCode)

	var events []responses.EventResponse
	err = json.Unmarshal(respBody, &events)
	s.NoError(err)

	s.Len(events, 0)
}

func (s *httpTestSuite) TestCheckDeleteEventNotFound() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf(
			"http://%s:%s/events?id=%s",
			s.cfg.HTTP.Host,
			s.cfg.HTTP.Port,
			s.eventID.String(),
		),
		nil)

	s.NoError(err)

	client := http.Client{}

	resp, err := client.Do(req)
	s.NoError(err)

	defer func() {
		_ = resp.Body.Close()
	}()

	s.Equal(http.StatusNotFound, resp.StatusCode)
}
