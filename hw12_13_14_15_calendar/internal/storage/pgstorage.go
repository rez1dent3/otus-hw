package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
	"time"
)

type PgStorage struct {
	dsn string
	db  *sqlx.DB
}

func NewPgStorage(dsn string) *PgStorage {
	return &PgStorage{dsn: dsn}
}

func (s *PgStorage) Connect(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, "postgres", s.dsn)
	if err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *PgStorage) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *PgStorage) CreateEvent(event Event) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Microsecond)
	defer cancel()

	_, err := s.db.NamedExecContext(
		ctx,
		`INSERT INTO events 
    			(id, title, description, start_at, end_at, user_id, created_at, updated_at) 
    		VALUES 
    			(:id, :title, :description, :start_at, :end_at, :user_id, NOW(), NOW())`, event)

	if err != nil {
		return false
	}

	return true
}

func (s *PgStorage) UpdateEvent(eventID uuid.UUID, event Event) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Microsecond)
	defer cancel()

	_, err := s.db.NamedExecContext(
		ctx,
		`UPDATE events 
    		SET
    		    title=:title,
    		    description=:description,
    		    start_at=:start_at,
    		    end_at=:end_at,
    		    updated_at=NOW()
    		WHERE
    		    id=:id`, event)

	if err != nil {
		return false
	}

	return true
}

func (s *PgStorage) DeleteEvent(eventID uuid.UUID) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Microsecond)
	defer cancel()

	execContext, err := s.db.ExecContext(ctx, `DELETE FROM events WHERE id=?`, eventID)
	if err != nil {
		return false
	}

	affected, err := execContext.RowsAffected()
	if err != nil {
		return false
	}

	return affected > 0
}

func (s *PgStorage) ListEventsDay(userID uuid.UUID, date time.Time) map[uuid.UUID]Event {
	panic("ListEventsDay")
}

func (s *PgStorage) ListEventsWeek(userID uuid.UUID, date time.Time) map[uuid.UUID]Event {
	panic("ListEventsWeek")
}

func (s *PgStorage) ListEventsMonth(userID uuid.UUID, date time.Time) map[uuid.UUID]Event {
	panic("ListEventsMonth")
}
