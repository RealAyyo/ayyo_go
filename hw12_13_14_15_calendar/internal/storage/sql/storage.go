package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/config"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/pkg/utils"
	"github.com/jackc/pgx/v5"
)

var (
	ErrEventNotFound = fmt.Errorf("event not found")
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
	err := s.db.QueryRow(
		ctx,
		"INSERT INTO events (title, date, duration, user_id) VALUES ($1, $2, $3, $4) RETURNING id",
		event.Title, event.Date, event.Duration, event.UserID,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, updated *storage.Event) (int, error) {
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

	query = fmt.Sprintf("%s WHERE id = $%d AND user_id = $%d RETURNING id", query[:len(query)-1], argsCount+1, argsCount+2)
	args = append(args, updated.ID, updated.UserID)

	err := s.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
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
		"SELECT id, title, date, duration::text, user_id FROM events WHERE user_id = $1 AND date >= $2 AND date <= $3",
		userID,
		dateFrom.Format("2006-01-02 15:04:05"),
		dateTo.Format("2006-01-02 15:04:05"),
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

func (s *Storage) CheckEventOverlaps(ctx context.Context, userID int, date time.Time, duration string) (bool, error) {
	durationParsed, err := utils.ParseDuration(duration)
	if err != nil {
		return false, err
	}
	endTime := date.Add(durationParsed)
	fmt.Println(date)
	fmt.Println(endTime)
	var count int
	err = s.db.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM events WHERE ((date <= $1 AND $1 < (date + duration)) OR ($2 <= date AND date < $2)) AND user_id = $3",
		date, endTime, userID,
	).Scan(&count)
	fmt.Println(count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
