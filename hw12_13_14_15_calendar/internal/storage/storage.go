package storage

import "context"

type ConnInterface interface {
	Connect(context.Context) error
	Close() error
}

func NewConnect(ctx context.Context, driver, dsn string) (ConnInterface, error) {
	if driver == "postgres" {
		pg := NewPgStorage(dsn)
		if err := pg.Connect(ctx); err != nil {
			return nil, err
		}

		return pg, nil
	}

	return NewMemStorage(), nil
}
