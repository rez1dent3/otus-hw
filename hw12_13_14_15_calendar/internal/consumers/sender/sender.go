package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/exchanges"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/logger"
)

type Sender struct {
	queue     exchanges.QueueInterface
	log       *logger.Logger
	queueName string
	storage   storage.SenderStorage
}

func New(log *logger.Logger, queue exchanges.QueueInterface, queueName string, storage storage.SenderStorage) *Sender {
	return &Sender{
		log:       log,
		queue:     queue,
		queueName: queueName,
		storage:   storage,
	}
}

func (s Sender) consume(body []byte) {
	notify := storage.Notify{}
	err := json.Unmarshal(body, &notify)
	if err != nil {
		s.log.Error("unmarshal body error: " + err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := s.storage.MarkAsDispatched(ctx, notify.EventID); err != nil {
		return
	}

	if notify.Description == "" {
		s.log.Info(fmt.Sprintf(
			"Hello. Event coming soon \"%s\", don't miss it. Start at: %s",
			notify.Title,
			notify.StartAt.String()))

		return
	}

	s.log.Info(fmt.Sprintf("Hello. Event coming soon \"%s\", don't miss it. Start at: %s\n%s",
		notify.Title,
		notify.StartAt.String(),
		notify.Description))
}

func (s Sender) Run(ctx context.Context) error {
	err := s.queue.Connect(ctx)
	if err != nil {
		return fmt.Errorf("error connect to rmq: %w", err)
	}

	err = s.queue.Receive(ctx, s.queueName, s.consume)
	if err != nil {
		return fmt.Errorf("error of send: %w", err)
	}

	return nil
}
