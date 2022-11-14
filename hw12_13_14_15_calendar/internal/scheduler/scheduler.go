package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/exchanges"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

type Scheduler struct {
	queue     exchanges.QueueInterface
	logger    *logger.Logger
	storage   storage.SchedulerStorage
	queueName string

	notifyDuration  time.Duration
	cleanerDuration time.Duration
}

func New(
	storage storage.SchedulerStorage,
	queue exchanges.QueueInterface,
	queueName string,
	logger *logger.Logger,
	notifyDuration time.Duration,
	cleanerDuration time.Duration,
) *Scheduler {
	return &Scheduler{
		queue:           queue,
		logger:          logger,
		storage:         storage,
		queueName:       queueName,
		notifyDuration:  notifyDuration,
		cleanerDuration: cleanerDuration,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	if err := s.queue.Connect(ctx); err != nil {
		return fmt.Errorf("error queue start: %w", err)
	}

	notifyTimer := time.NewTicker(s.notifyDuration)
	removeTimer := time.NewTicker(s.cleanerDuration)

	s.logger.Info("Calendar scheduler started...")

	defer removeTimer.Stop()
	defer notifyTimer.Stop()

	for {
		select {
		case <-notifyTimer.C:
			err := s.notify(ctx)
			if err != nil {
				s.logger.Error("notification error: " + err.Error())
			}
		case <-removeTimer.C:
			if err := s.remove(ctx); err != nil {
				s.logger.Error("remove events error: " + err.Error())
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *Scheduler) notify(ctx context.Context) error {
	notifies, err := s.storage.ListToSendNotifies(ctx, time.Now())
	if err != nil {
		return err
	}

	s.logger.Info("found notifies to send: " + strconv.Itoa(len(notifies)))
	for _, notify := range notifies {
		body, err := json.Marshal(notify)
		if err != nil {
			s.logger.Error("marshal error: " + err.Error())
			continue
		}

		s.logger.Info(fmt.Sprintf("sending an notify %s to the queue", notify.EventID.String()))
		if err := s.queue.Sent(ctx, s.queueName, body); err != nil {
			s.logger.Error("sent error: " + err.Error())
			continue
		}

		if err := s.storage.MarkAsSent(ctx, notify.EventID); err != nil {
			s.logger.Error("set is notified error: " + err.Error())
		}
	}

	return nil
}

func (s *Scheduler) remove(ctx context.Context) error {
	return s.storage.RemoveOldEvents(ctx, time.Now().AddDate(-1, 0, 0))
}
