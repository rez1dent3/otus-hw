package rabbitmq

import (
	"context"
	"fmt"
	"net"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/exchanges"
)

type impl struct {
	conn *amqp.Connection

	dsn string
}

func New(dsn string) exchanges.QueueInterface {
	return &impl{dsn: dsn}
}

func (i *impl) Connect(ctx context.Context) error {
	conn, err := amqp.DialConfig(i.dsn, amqp.Config{
		Dial: func(network, addr string) (net.Conn, error) {
			dialer := net.Dialer{}

			return dialer.DialContext(ctx, network, addr)
		},
	})

	if err != nil {
		return err
	}

	i.conn = conn

	return nil
}

func (i *impl) getQueue(name string, ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
}

func (i *impl) Close() error {
	if i.conn.IsClosed() {
		return nil
	}

	return i.conn.Close()
}

func (i *impl) Receive(ctx context.Context, name string, callback func(body []byte)) error {
	ch, err := i.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed open channel: %w", err)
	}

	q, err := i.getQueue(name, ch)
	if err != nil {
		return fmt.Errorf("failed queue declare: %w", err)
	}

	msg, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed open channel: %w", err)
	}

	go func() {
		for d := range msg {
			callback(d.Body)
		}
	}()

	<-ctx.Done()

	return nil
}

func (i *impl) Sent(ctx context.Context, name string, body []byte) error {
	ch, err := i.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed open channel: %w", err)
	}

	q, err := i.getQueue(name, ch)
	if err != nil {
		return fmt.Errorf("failed queue declare: %w", err)
	}

	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		return fmt.Errorf("failed publish: %w", err)
	}

	return ch.Close()
}
