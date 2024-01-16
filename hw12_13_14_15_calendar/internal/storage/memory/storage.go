package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/app"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrEventNotFound = errors.New("event for update not found")
)

type EventsMap map[int]map[int]*storage.Event

type Storage struct {
	count  int
	events EventsMap
	mu     sync.RWMutex //nolint:unused
}

func (s *Storage) AddEvent(ctx context.Context, event *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if event.UserId == 0 {
		return app.ErrUserIdRequired
	}

	if event.Title == "" {
		return app.ErrTitleRequired
	}

	if event.Date.IsZero() {
		return app.ErrDateRequired
	}

	if event.Duration == "" {
		return app.ErrDurationRequired
	}

	id := s.count

	if s.events[event.UserId] == nil {
		s.events[event.UserId] = make(map[int]*storage.Event)
	}

	s.events[event.UserId][id] = event
	s.count += 1

	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, updated *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if updated.UserId == 0 {
		return app.ErrUserIdRequired
	}

	if updated.ID == 0 {
		return app.ErrEventIdRequired
	}

	findEvent, ok := s.events[updated.UserId][updated.ID]
	if !ok {
		return ErrEventNotFound
	}

	if updated.Title != "" {
		findEvent.Title = updated.Title
	}

	if updated.Duration != "" {
		findEvent.Duration = updated.Duration
	}

	if !updated.Date.IsZero() {
		findEvent.Date = updated.Date
	}

	s.events[updated.UserId][updated.ID] = findEvent
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id int, userId int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[userId][id]; !ok {
		return ErrEventNotFound
	}

	delete(s.events[userId], id)

	return nil
}

func (s *Storage) ListEvents(ctx context.Context, userId int, dateFrom time.Time, dateTo time.Time) ([]storage.Event, error) {
	var results []storage.Event

	for id, event := range s.events[userId] {
		if (event.Date.After(dateFrom) || event.Date.Equal(dateFrom)) && (event.Date.Before(dateTo) || event.Date.Equal(dateTo)) {
			results = append(results, storage.Event{
				ID:       id,
				Title:    event.Title,
				Date:     event.Date,
				Duration: event.Duration,
				UserId:   event.UserId,
			})
		}
	}
	return results, nil
}

func New() (*Storage, error) {
	return &Storage{
		count:  1,
		events: make(EventsMap),
	}, nil
}
