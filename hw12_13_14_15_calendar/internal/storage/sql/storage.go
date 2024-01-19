package sqlstorage

import (
	"context"
	"time"

	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/app"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/config"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage"
	"github.com/jackc/pgx/v5"
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

func (s *Storage) AddEvent(ctx context.Context, event *storage.Event) error {
	_, err := s.db.Exec(
		ctx,
		"INSERT INTO events (title, date, duration, user_id) VALUES ($1, $2, $3, $4)",
		event.Title, event.Date, event.Duration, event.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, updated *storage.Event) error {
	if updated.UserID == 0 {
		return app.ErrUserIDRequired
	}

	_, err := s.db.Exec(ctx, "UPDATE events SET title = $1, duration = $2, date = $3, WHERE id = $5 AND user_id = $6", updated.Title, updated.Duration, updated.Date, updated.ID, updated.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id int, userID int) error {
	_, err := s.db.Exec(
		ctx,
		"DELETE FROM events WHERE id = $1 AND user_id = $2",
		id, userID,
	)
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
		"SELECT id, title, date, duration::text, user_id FROM events WHERE user_id = $1 AND date >= $2 AND date <= $3",
		userID,
		dateFrom,
		dateTo,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var event storage.Event

		err = rows.Scan(&event.ID, &event.Title, &event.Date, &event.Duration, &event.UserID)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}
