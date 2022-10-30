package storage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rez1dent3/otus-hw/hw12_13_14_15_calendar/pkg/uuid"
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
	return s.db.Close()
}

func (s *PgStorage) CreateEvent(ctx context.Context, event Event) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rows, err := s.db.NamedExecContext(
		ctx,
		`INSERT INTO events 
    			(id, title, description, start_at, end_at, user_id, created_at, updated_at) 
    		VALUES 
    			(:id, :title, :description, :start_at, :end_at, :user_id, NOW(), NOW());`,
		event)
	if err != nil {
		return false
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		return false
	}

	return affected > 0
}

func (s *PgStorage) UpdateEvent(ctx context.Context, eventID uuid.UUID, event Event) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rows, err := s.db.ExecContext(
		ctx,
		`UPDATE events 
    		SET
    		    title=$1,
    		    description=$2,
    		    start_at=$3,
    		    end_at=$4,
    		    remind_for=$5,
    		    user_id=$6,
    		    updated_at=NOW()
    		WHERE
    		    id=$7`,
		event.Title,
		event.Description,
		event.StartAt,
		event.EndAt,
		event.RemindFor,
		event.UserID,
		eventID,
	)
	if err != nil {
		return false
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		return false
	}

	return affected > 0
}

func (s *PgStorage) DeleteEvent(ctx context.Context, eventID uuid.UUID) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	execContext, err := s.db.ExecContext(ctx, `DELETE FROM events WHERE id=$1`, eventID)
	if err != nil {
		return false
	}

	affected, err := execContext.RowsAffected()
	if err != nil {
		return false
	}

	return affected > 0
}

func (s *PgStorage) ListEventsDay(ctx context.Context, userID uuid.UUID, date time.Time) map[uuid.UUID]Event {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rows, err := s.db.QueryxContext(
		ctx,
		`SELECT *
			FROM events
			WHERE user_id=$1 AND
        		date_trunc('day', $2::date)::date BETWEEN 
        			date_trunc('day', start_at)::date AND end_at::date`, userID, date)
	if err != nil {
		return nil
	}

	defer func() {
		rows.Close()
	}()

	result := make(map[uuid.UUID]Event)
	event := Event{}
	for rows.Next() {
		err := rows.StructScan(&event)
		if err != nil {
			return nil
		}

		event := event
		result[event.ID] = event
	}

	return result
}

func (s *PgStorage) ListEventsWeek(ctx context.Context, userID uuid.UUID, date time.Time) map[uuid.UUID]Event {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rows, err := s.db.QueryxContext(
		ctx,
		`SELECT *
			FROM events
			WHERE user_id=$1 AND
        		date_trunc('week', $2::date)::date BETWEEN 
        			date_trunc('week', start_at)::date AND end_at::date`, userID, date)
	if err != nil {
		return nil
	}

	defer func() {
		rows.Close()
	}()

	result := make(map[uuid.UUID]Event)
	event := Event{}
	for rows.Next() {
		err := rows.StructScan(&event)
		if err != nil {
			return nil
		}

		event := event
		result[event.ID] = event
	}

	return result
}

func (s *PgStorage) ListEventsMonth(ctx context.Context, userID uuid.UUID, date time.Time) map[uuid.UUID]Event {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rows, err := s.db.QueryxContext(
		ctx,
		`SELECT *
			FROM events
			WHERE user_id=$1 AND
        		date_trunc('month', $2::date)::date BETWEEN 
        			date_trunc('month', start_at)::date AND end_at::date`, userID, date)
	if err != nil {
		return nil
	}

	defer func() {
		rows.Close()
	}()

	result := make(map[uuid.UUID]Event)
	event := Event{}
	for rows.Next() {
		err := rows.StructScan(&event)
		if err != nil {
			return nil
		}

		event := event
		result[event.ID] = event
	}

	return result
}
