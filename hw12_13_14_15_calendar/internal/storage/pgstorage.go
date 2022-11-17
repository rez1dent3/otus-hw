package storage

import (
	"context"
	"strings"
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

func (s *PgStorage) CreateEvent(parentCtx context.Context, event Event) (bool, error) {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second)
	defer cancel()

	rows, err := s.db.NamedExecContext(
		ctx,
		`INSERT INTO events 
    			(id, title, description, start_at, end_at, user_id, in_queue, is_dispatched, remind_for, created_at, updated_at) 
    		VALUES 
    			(:id, :title, :description, :start_at, :end_at, :user_id, :in_queue, :is_dispatched, :remind_for, NOW(), NOW());`,
		event)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return false, ErrUnableDuplicate
		}

		return false, err
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		return false, err
	}

	if affected > 0 {
		return true, nil
	}

	return false, ErrNotFound
}

func (s *PgStorage) UpdateEvent(parentCtx context.Context, eventID uuid.UUID, event Event) (bool, error) {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second)
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
		return false, err
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		return false, err
	}

	if affected > 0 {
		return true, nil
	}

	return false, ErrNotFound
}

func (s *PgStorage) DeleteEvent(parentCtx context.Context, eventID uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second)
	defer cancel()

	execContext, err := s.db.ExecContext(ctx, `DELETE FROM events WHERE id=$1`, eventID)
	if err != nil {
		return false, err
	}

	affected, err := execContext.RowsAffected()
	if err != nil {
		return false, err
	}

	if affected > 0 {
		return true, nil
	}

	return false, ErrNotFound
}

func (s *PgStorage) ListEventsDay(parentCtx context.Context, userID uuid.UUID, date time.Time) map[uuid.UUID]Event {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second)
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

func (s *PgStorage) ListEventsWeek(parentCtx context.Context, userID uuid.UUID, date time.Time) map[uuid.UUID]Event {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second)
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

func (s *PgStorage) ListEventsMonth(parentCtx context.Context, userID uuid.UUID, date time.Time) map[uuid.UUID]Event {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second)
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

func (s *PgStorage) ListToSendNotifies(parentCtx context.Context, date time.Time) ([]Notify, error) {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second)
	defer cancel()

	rows, err := s.db.QueryxContext(
		ctx,
		`SELECT id,
       				title,
       				description,
       				start_at,
       				end_at,
       				user_id
			FROM events
			WHERE in_queue=false 
			  AND (
			      remind_for IS NOT NULL AND start_at-remind_for<=$1
			      OR
			      remind_for IS NULL AND start_at<=$2
			  )
		  	LIMIT 1000`, date, date)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	var result []Notify
	notify := Notify{}
	for rows.Next() {
		err := rows.StructScan(&notify)
		if err != nil {
			return nil, err
		}

		notify := notify
		result = append(result, notify)
	}

	return result, nil
}

func (s *PgStorage) RemoveOldEvents(parentCtx context.Context, date time.Time) error {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		`DELETE FROM events
    		WHERE end_at<=$1`,
		date,
	)

	return err
}

func (s *PgStorage) MarkInQueue(parentCtx context.Context, eventID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		`UPDATE events 
    		SET
    		    in_queue=true
    		WHERE
    		    id=$1`,
		eventID,
	)

	return err
}

func (s *PgStorage) MarkAsDispatched(parentCtx context.Context, eventID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		`UPDATE events 
    		SET
    		    is_dispatched=true
    		WHERE
    		    id=$1`,
		eventID,
	)

	return err
}
