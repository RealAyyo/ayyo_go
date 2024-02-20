package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/config"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/pkg/utils"
	"github.com/jackc/pgx/v5"
)

var (
	ErrEventNotFound = fmt.Errorf("event not found")
	ErrDateBusy      = errors.New("date is busy")
)

type Closer interface {
	Close(ctx context.Context) error
}

type Storage struct {
	db *pgx.Conn
}

func New(ctx context.Context, conf config.DBConf) (*Storage, error) {
	sqlStorage := &Storage{}

	err := sqlStorage.Connect(ctx, conf)
	if err != nil {
		return nil, err
	}

	return sqlStorage, nil
}

func (s *Storage) Connect(ctx context.Context, conf config.DBConf) error {
	connString := "postgres://" + conf.Username + ":" + conf.Password + "@" + conf.Host + ":" + conf.Port + "/" + conf.Database

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return err
	}
	s.db = conn

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	if s.db.IsClosed() {
		return nil
	}

	err := s.db.Close(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) AddEvent(ctx context.Context, event *storage.Event) (int, error) {
	var id int

	var notificationTime *string
	if event.NotificationTime != "" {
		notificationTime = &event.NotificationTime
	}

	err := s.db.QueryRow(
		ctx,
		"INSERT INTO events (title, date, duration, user_id, description, notification_time) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		event.Title, event.Date, event.Duration, event.UserID, event.Description, notificationTime,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, updated *storage.Event) error {
	var id int
	args := []interface{}{}
	argsCount := 0
	query := "UPDATE events SET"

	if updated.Title != "" {
		argsCount++
		query += fmt.Sprintf(" title = $%d,", argsCount)
		args = append(args, updated.Title)
	}
	if !updated.Date.IsZero() {
		argsCount++
		query += fmt.Sprintf(" date = $%d,", argsCount)
		args = append(args, updated.Date)
	}
	if updated.Duration != "" {
		argsCount++
		query += fmt.Sprintf(" duration = $%d,", argsCount)
		args = append(args, updated.Duration)
	}
	if updated.Description != "" {
		argsCount++
		query += fmt.Sprintf(" description = $%d,", argsCount)
		args = append(args, updated.Duration)
	}
	if updated.NotificationTime != "" {
		argsCount++
		query += fmt.Sprintf(" notification_time = $%d,", argsCount)
		args = append(args, updated.Duration)
	}

	query = fmt.Sprintf("%s WHERE id = $%d AND user_id = $%d RETURNING id", query[:len(query)-1], argsCount+1, argsCount+2)
	args = append(args, updated.ID, updated.UserID)

	err := s.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id int, userID int) error {
	r, err := s.db.Exec(
		ctx,
		"DELETE FROM events WHERE id = $1 AND user_id = $2",
		id, userID,
	)
	if r.RowsAffected() == 0 {
		return ErrEventNotFound
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) ListEvents(
	ctx context.Context, userID int, dateFrom time.Time, dateTo time.Time,
) ([]storage.Event, error) {
	var events []storage.Event

	rows, err := s.db.Query(
		ctx,
		"SELECT id, title, date, duration::text, user_id, description, notification_time::text FROM events WHERE user_id = $1 AND date >= $2 AND date <= $3",
		userID,
		dateFrom.Format(time.RFC3339),
		dateTo.Format(time.RFC3339),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var event storage.Event
		var description sql.NullString
		var notificationTime sql.NullString

		err = rows.Scan(&event.ID, &event.Title, &event.Date, &event.Duration, &event.UserID, &description, &notificationTime)
		if err != nil {
			return nil, err
		}

		if description.Valid {
			event.Description = description.String
		}
		if notificationTime.Valid {
			event.NotificationTime = notificationTime.String
		}

		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) EventsCleanUp(ctx context.Context) error {
	yearAgo := time.Now().AddDate(-1, 0, 0)
	_, err := s.db.Exec(
		ctx,
		"DELETE FROM events WHERE date < $1",
		yearAgo,
	)
	return err
}

func (s *Storage) GetEventsToNotify(ctx context.Context) ([]storage.Event, error) {
	var events []storage.Event

	rows, err := s.db.Query(
		ctx,
		`SELECT id, title, date, duration::text, user_id, description, notification_time::text 
   FROM events 
   WHERE date - COALESCE(notification_time::interval, '0s'::interval) <= NOW()`,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var event storage.Event
		var description sql.NullString

		err = rows.Scan(&event.ID, &event.Title, &event.Date, &event.Duration, &event.UserID, &description, &event.NotificationTime)
		if err != nil {
			return nil, err
		}

		if description.Valid {
			event.Description = description.String
		}

		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) CheckEventOverlaps(ctx context.Context, userID int, date time.Time, duration string) error {
	durationParsed, err := utils.ParseDuration(duration)
	if err != nil {
		return err
	}
	endTime := date.Add(durationParsed)

	var exists bool
	err = s.db.QueryRow(
		ctx,
		"SELECT EXISTS (SELECT 1 FROM events WHERE ((date <= $1 AND $1 < (date + duration)) OR ($2 <= date AND date < $2)) AND user_id = $3)",
		date, endTime, userID,
	).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return ErrDateBusy
	}
	return nil
}
