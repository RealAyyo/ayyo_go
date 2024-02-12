package app

import (
	"context"
	"errors"
	"time"

	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrDateBusy         = errors.New("invalid string")
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
	AddEvent(ctx context.Context, event *storage.Event) error
	UpdateEvent(ctx context.Context, updated *storage.Event) error
	DeleteEvent(ctx context.Context, id int, userID int) error
	ListEvents(ctx context.Context, userID int, dateFrom time.Time, dateTo time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage StorageService) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event *storage.Event) error {
	return a.storage.AddEvent(ctx, event)
}

func (a *App) GetEventsForRange(
	ctx context.Context, userID int, dateFrom time.Time, dateRange int,
) ([]storage.Event, error) {
	var dateTo time.Time
	switch dateRange {
	case DAY:
		dateTo = dateFrom.AddDate(0, 0, 1)
	case WEEK:
		dateTo = dateFrom.AddDate(0, 0, 7)
	case MONTH:
		dateTo = dateFrom.AddDate(0, 1, 0)
	default:
		return nil, ErrDateRange
	}

	listEvents, err := a.storage.ListEvents(ctx, userID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}

	return listEvents, nil
}
