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
	s.False(events[0].IsDispatched)
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
	s.False(events[0].IsDispatched)
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

func (s *httpTestSuite) TestCheckDayListEvents() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client := http.Client{}

	userID := uuid.Gen()
	now := time.Date(2022, 11, 21, 0, 0, 0, 0, time.UTC)
	var eventRequests []requests.EventRequest
	for i := 0; i < 3; i++ {
		// 21, 28, 5
		eventID := uuid.Gen()
		event := requests.EventRequest{
			ID:        eventID.String(),
			StartAt:   now.Add(time.Duration(i*7) * 24 * time.Hour),
			EndAt:     now.Add(time.Duration(i*7+1) * 24 * time.Hour),
			UserID:    userID.String(),
			RemindFor: nil,
		}
		eventRequests = append(eventRequests, event)

		reqBody, err := json.Marshal(event)
		s.NoError(err)

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			fmt.Sprintf("http://%s:%s/events", s.cfg.HTTP.Host, s.cfg.HTTP.Port),
			bytes.NewReader(reqBody))

		s.NoError(err)

		resp, err := client.Do(req)
		s.NoError(err)

		s.Equal(http.StatusOK, resp.StatusCode)
		s.NoError(resp.Body.Close())
	}

	// day
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"http://%s:%s/events?type=ListEventsDay&userId=%s&timestamp=%d",
			s.cfg.HTTP.Host,
			s.cfg.HTTP.Port,
			userID.String(),
			now.Add(24*time.Hour).Unix(),
		),
		nil)

	resp, err := client.Do(req)
	s.NoError(err)

	respBody, err := io.ReadAll(resp.Body)
	s.NoError(resp.Body.Close())
	s.Equal(http.StatusOK, resp.StatusCode)

	var eventResp1 []responses.EventResponse

	err = json.Unmarshal(respBody, &eventResp1)
	s.NoError(err)
	s.Len(eventResp1, 1)
	s.Equal(eventRequests[0].ID, eventResp1[0].ID)

	// not found
	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"http://%s:%s/events?type=ListEventsDay&userId=%s&timestamp=%d",
			s.cfg.HTTP.Host,
			s.cfg.HTTP.Port,
			userID.String(),
			now.Add(2*24*time.Hour).Unix(),
		),
		nil)

	resp, err = client.Do(req)
	s.NoError(err)

	respBody, err = io.ReadAll(resp.Body)
	s.NoError(resp.Body.Close())
	s.Equal(http.StatusOK, resp.StatusCode)

	var eventResp2 []responses.EventResponse

	err = json.Unmarshal(respBody, &eventResp2)
	s.NoError(err)
	s.Len(eventResp2, 0)
}

func (s *httpTestSuite) TestCheckWeekListEvents() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client := http.Client{}

	userID := uuid.Gen()
	now := time.Date(2022, 11, 21, 0, 0, 0, 0, time.UTC)
	var eventRequests []requests.EventRequest
	for i := 0; i < 3; i++ {
		// 21, 28, 5
		eventID := uuid.Gen()
		event := requests.EventRequest{
			ID:        eventID.String(),
			StartAt:   now.Add(time.Duration(i*7) * 24 * time.Hour),
			EndAt:     now.Add(time.Duration(i*7+1) * 24 * time.Hour),
			UserID:    userID.String(),
			RemindFor: nil,
		}
		eventRequests = append(eventRequests, event)

		reqBody, err := json.Marshal(event)
		s.NoError(err)

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			fmt.Sprintf("http://%s:%s/events", s.cfg.HTTP.Host, s.cfg.HTTP.Port),
			bytes.NewReader(reqBody))

		s.NoError(err)

		resp, err := client.Do(req)
		s.NoError(err)

		s.Equal(http.StatusOK, resp.StatusCode)
		s.NoError(resp.Body.Close())
	}

	// week
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"http://%s:%s/events?type=ListEventsWeek&userId=%s&timestamp=%d",
			s.cfg.HTTP.Host,
			s.cfg.HTTP.Port,
			userID.String(),
			now.Add(24*time.Hour).Unix(),
		),
		nil)

	resp, err := client.Do(req)
	s.NoError(err)

	respBody, err := io.ReadAll(resp.Body)
	s.NoError(resp.Body.Close())
	s.Equal(http.StatusOK, resp.StatusCode)

	var eventResp1 []responses.EventResponse

	err = json.Unmarshal(respBody, &eventResp1)
	s.NoError(err)
	s.Len(eventResp1, 1)
	s.Equal(eventRequests[0].ID, eventResp1[0].ID)

	// next week
	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"http://%s:%s/events?type=ListEventsWeek&userId=%s&timestamp=%d",
			s.cfg.HTTP.Host,
			s.cfg.HTTP.Port,
			userID.String(),
			now.Add(7*24*time.Hour).Unix(),
		),
		nil)

	resp, err = client.Do(req)
	s.NoError(err)

	respBody, err = io.ReadAll(resp.Body)
	s.NoError(resp.Body.Close())
	s.Equal(http.StatusOK, resp.StatusCode)

	var eventResp2 []responses.EventResponse

	err = json.Unmarshal(respBody, &eventResp2)
	s.NoError(err)
	s.Len(eventResp2, 1)
	s.Equal(eventRequests[1].ID, eventResp2[0].ID)

	// next week
	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"http://%s:%s/events?type=ListEventsWeek&userId=%s&timestamp=%d",
			s.cfg.HTTP.Host,
			s.cfg.HTTP.Port,
			userID.String(),
			now.Add(14*24*time.Hour).Unix(),
		),
		nil)

	resp, err = client.Do(req)
	s.NoError(err)

	respBody, err = io.ReadAll(resp.Body)
	s.NoError(resp.Body.Close())
	s.Equal(http.StatusOK, resp.StatusCode)

	var eventResp3 []responses.EventResponse

	err = json.Unmarshal(respBody, &eventResp3)
	s.NoError(err)
	s.Len(eventResp3, 1)
	s.Equal(eventRequests[2].ID, eventResp3[0].ID)

	// not found - week
	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"http://%s:%s/events?type=ListEventsWeek&userId=%s&timestamp=%d",
			s.cfg.HTTP.Host,
			s.cfg.HTTP.Port,
			userID.String(),
			now.Add(21*24*time.Hour).Unix(),
		),
		nil)

	resp, err = client.Do(req)
	s.NoError(err)

	respBody, err = io.ReadAll(resp.Body)
	s.NoError(resp.Body.Close())
	s.Equal(http.StatusOK, resp.StatusCode)

	var eventResp4 []responses.EventResponse

	err = json.Unmarshal(respBody, &eventResp4)
	s.NoError(err)
	s.Len(eventResp4, 0)
}

func (s *httpTestSuite) TestCheckMonthListEvents() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client := http.Client{}

	userID := uuid.Gen()
	now := time.Date(2022, 11, 21, 0, 0, 0, 0, time.UTC)
	var eventRequests []requests.EventRequest
	for i := 0; i < 3; i++ {
		// 21, 28, 5
		eventID := uuid.Gen()
		event := requests.EventRequest{
			ID:        eventID.String(),
			StartAt:   now.Add(time.Duration(i*7) * 24 * time.Hour),
			EndAt:     now.Add(time.Duration(i*7+1) * 24 * time.Hour),
			UserID:    userID.String(),
			RemindFor: nil,
		}
		eventRequests = append(eventRequests, event)

		reqBody, err := json.Marshal(event)
		s.NoError(err)

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			fmt.Sprintf("http://%s:%s/events", s.cfg.HTTP.Host, s.cfg.HTTP.Port),
			bytes.NewReader(reqBody))

		s.NoError(err)

		resp, err := client.Do(req)
		s.NoError(err)

		s.Equal(http.StatusOK, resp.StatusCode)
		s.NoError(resp.Body.Close())
	}

	// month
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"http://%s:%s/events?type=ListEventsMonth&userId=%s&timestamp=%d",
			s.cfg.HTTP.Host,
			s.cfg.HTTP.Port,
			userID.String(),
			now.Unix(),
		),
		nil)

	resp, err := client.Do(req)
	s.NoError(err)

	respBody, err := io.ReadAll(resp.Body)
	s.NoError(resp.Body.Close())
	s.Equal(http.StatusOK, resp.StatusCode)

	var eventResp1 []responses.EventResponse

	err = json.Unmarshal(respBody, &eventResp1)
	s.NoError(err)
	s.Len(eventResp1, 2)
	s.Equal(eventRequests[0].ID, eventResp1[0].ID)
	s.Equal(eventRequests[1].ID, eventResp1[1].ID)

	// next month
	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"http://%s:%s/events?type=ListEventsMonth&userId=%s&timestamp=%d",
			s.cfg.HTTP.Host,
			s.cfg.HTTP.Port,
			userID.String(),
			now.AddDate(0, 1, 0).Unix(),
		),
		nil)

	resp, err = client.Do(req)
	s.NoError(err)

	respBody, err = io.ReadAll(resp.Body)
	s.NoError(resp.Body.Close())
	s.Equal(http.StatusOK, resp.StatusCode)

	var eventResp2 []responses.EventResponse

	err = json.Unmarshal(respBody, &eventResp2)
	s.NoError(err)
	s.Len(eventResp2, 1)
	s.Equal(eventRequests[2].ID, eventResp2[0].ID)

	// not found - month
	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"http://%s:%s/events?type=ListEventsMonth&userId=%s&timestamp=%d",
			s.cfg.HTTP.Host,
			s.cfg.HTTP.Port,
			userID.String(),
			now.AddDate(0, 2, 0).Unix(),
		),
		nil)

	resp, err = client.Do(req)
	s.NoError(err)

	respBody, err = io.ReadAll(resp.Body)
	s.NoError(resp.Body.Close())
	s.Equal(http.StatusOK, resp.StatusCode)

	var eventResp3 []responses.EventResponse

	err = json.Unmarshal(respBody, &eventResp3)
	s.NoError(err)
	s.Len(eventResp3, 0)
}

func (s *httpTestSuite) TestCheckSenderApp() {
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

	var events1 []responses.EventResponse
	err = json.Unmarshal(respBody, &events1)
	s.NoError(err)

	s.Equal(s.event.ID, events1[0].ID)
	s.False(events1[0].IsDispatched)

	// wait until the scheduler and sender work
	time.Sleep(5 * time.Second)

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second)
	defer cancel2()

	req, err = http.NewRequestWithContext(
		ctx2,
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

	respBody, err = io.ReadAll(resp.Body)
	s.NoError(resp.Body.Close())
	s.Equal(http.StatusOK, resp.StatusCode)

	var events2 []responses.EventResponse
	err = json.Unmarshal(respBody, &events2)
	s.NoError(err)

	s.Equal(s.event.ID, events2[0].ID)
	s.True(events2[0].IsDispatched)
}