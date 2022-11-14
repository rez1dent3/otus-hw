package exchanges

import "context"

type QueueInterface interface {
	Connect(ctx context.Context) error
	Receive(ctx context.Context, name string, callback func(body []byte)) error
	Sent(ctx context.Context, name string, body []byte) error
	Close() error
}
