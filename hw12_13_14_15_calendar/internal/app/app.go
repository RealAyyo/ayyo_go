package app

import (
	"context"
	"errors"
	"time"

	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrDateBusy         = errors.New("date is busy")
	ErrUserCantChange   = errors.New("user id can't change")
	ErrDateRange        = errors.New("invalid date range type")
	ErrEventIDRequired  = errors.New("event id required")
	ErrUserIDRequired   = errors.New("user id is required")
	ErrDateRequired     = errors.New("date is required")
	ErrDurationRequired = errors.New("duration is required")
	ErrTitleRequired    = errors.New("title is required")
)

const (
	DAY = iota
	WEEK
	MONTH
)

type App struct {
	logger  Logger
	storage StorageService
}

type Logger interface {
	Info(msg string, attrs ...any)
	Error(msg string, attrs ...any)
	Debug(msg string, attrs ...any)
	Warn(msg string, attrs ...any)
}

type StorageService interface {
	AddEvent(ctx context.Context, event *storage.Event) (int, error)
	UpdateEvent(ctx context.Context, updated *storage.Event) (int, error)
	DeleteEvent(ctx context.Context, id int, userID int) error
	ListEvents(ctx context.Context, userID int, dateFrom time.Time, dateTo time.Time) ([]storage.Event, error)
	CheckEventOverlaps(ctx context.Context, date time.Time, duration string) (bool, error)
}

func New(logger Logger, storage StorageService) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) UpdateEvent(ctx context.Context, event *storage.Event) (int, error) {
	if event.UserID == 0 {
		return 0, ErrUserIDRequired
	}
	if event.ID == 0 {
		return 0, ErrEventIDRequired
	}

	id, err := a.storage.UpdateEvent(ctx, event)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (a *App) DeleteEvent(ctx context.Context, eventID int, userID int) error {
	if userID == 0 {
		return ErrUserIDRequired
	}
	if eventID == 0 {
		return ErrEventIDRequired
	}

	return a.storage.DeleteEvent(ctx, eventID, userID)
}

func (a *App) CreateEvent(ctx context.Context, event *storage.Event) (int, error) {
	checkOverlap, err := a.storage.CheckEventOverlaps(ctx, event.Date, event.Duration)
	if err != nil {
		return 0, err
	}
	if checkOverlap {
		return 0, ErrDateBusy
	}

	return a.storage.AddEvent(ctx, event)
}

func (a *App) GetEventsForRange(ctx context.Context, userID int, dateFrom int64, dateTo int64) ([]storage.Event, error) {
	parsedDateFrom := time.Unix(dateFrom, 0)
	parsedDateTo := time.Unix(dateTo, 0)

	listEvents, err := a.storage.ListEvents(ctx, userID, parsedDateFrom, parsedDateTo)
	if err != nil {
		return nil, err
	}

	return listEvents, nil
}
