package sqlstorage

import (
	"context"
	"fmt"
	"strings"
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
	db       *pgx.Conn
	username string
	password string
	host     string
	port     string
	database string
}

func New(ctx context.Context, conf config.DbConf) (*Storage, error) {
	sqlStorage := &Storage{
		username: conf.Username,
		password: conf.Password,
		host:     conf.Host,
		port:     conf.Port,
		database: conf.Database,
	}

	err := sqlStorage.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return sqlStorage, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	connString := "postgres://" + s.username + ":" + s.password + "@" + s.host + ":" + s.port + "/" + s.database
	fmt.Println(connString)
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
		event.Title, event.Date, event.Duration, event.UserId)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, updated *storage.Event) error {
	var setParts []string
	var args []interface{}
	argIndex := 1

	if updated.UserId == 0 {
		return app.ErrUserIdRequired
	}

	setParts = append(setParts, fmt.Sprintf("user_id = $%d", argIndex))
	args = append(args, updated.UserId)
	argIndex++

	if updated.Title != "" {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, updated.Title)
		argIndex++
	}

	if updated.Duration != "" {
		setParts = append(setParts, fmt.Sprintf("duration = $%d", argIndex))
		args = append(args, updated.Duration)
		argIndex++
	}

	if updated.Date.IsZero() {
		setParts = append(setParts, fmt.Sprintf("date = $%d", argIndex))
		args = append(args, updated.Date)
		argIndex++
	}

	queryString := fmt.Sprintf("UPDATE events SET %s WHERE id = $%d", strings.Join(setParts, ", "), argIndex)
	args = append(args, updated.ID)

	_, err := s.db.Exec(ctx, queryString, args...)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id int, userId string) error {
	_, err := s.db.Exec(
		ctx,
		"DELETE FROM events WHERE id = $1 AND user_id = $2",
		id, userId,
	)

	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) ListEvents(ctx context.Context, userId int, dateFrom time.Time, dateTo time.Time) ([]storage.Event, error) {
	var events []storage.Event

	rows, err := s.db.Query(
		ctx,
		"SELECT id, title, date, duration::text, user_id FROM events WHERE user_id = $1 AND date >= $2 AND date <= $3",
		userId,
		dateFrom,
		dateTo,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var event storage.Event

		err = rows.Scan(&event.ID, &event.Title, &event.Date, &event.Duration, &event.UserId)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}
