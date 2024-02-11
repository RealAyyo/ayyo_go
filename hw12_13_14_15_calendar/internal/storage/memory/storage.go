package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/app"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage"
)

var ErrEventNotFound = errors.New("event for update not found")

type EventsMap map[int]map[int]*storage.Event

type Storage struct {
	count  int
	events EventsMap
	mu     sync.RWMutex
}

func (s *Storage) AddEvent(_ context.Context, event *storage.Event) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if event.UserID == 0 {
		return 0, app.ErrUserIDRequired
	}

	if event.Title == "" {
		return 0, app.ErrTitleRequired
	}

	if event.Date.IsZero() {
		return 0, app.ErrDateRequired
	}

	if event.Duration == "" {
		return 0, app.ErrDurationRequired
	}

	id := s.count

	if s.events[event.UserID] == nil {
		s.events[event.UserID] = make(map[int]*storage.Event)
	}

	s.events[event.UserID][id] = event
	s.count++

	return id, nil
}

func (s *Storage) UpdateEvent(_ context.Context, updated *storage.Event) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if updated.UserID == 0 {
		return 0, app.ErrUserIDRequired
	}

	if updated.ID == 0 {
		return 0, app.ErrEventIDRequired
	}

	findEvent, ok := s.events[updated.UserID][updated.ID]
	if !ok {
		return 0, ErrEventNotFound
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

	s.events[updated.UserID][updated.ID] = findEvent
	return updated.ID, nil
}

func (s *Storage) DeleteEvent(_ context.Context, id int, userID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[userID][id]; !ok {
		return ErrEventNotFound
	}

	delete(s.events[userID], id)

	return nil
}

func (s *Storage) ListEvents(
	_ context.Context, userID int, dateFrom time.Time, dateTo time.Time,
) ([]storage.Event, error) {
	var results []storage.Event

	for id, event := range s.events[userID] {
		if (event.Date.After(dateFrom) || event.Date.Equal(dateFrom)) &&
			(event.Date.Before(dateTo) || event.Date.Equal(dateTo)) {
			results = append(results, storage.Event{
				ID:       id,
				Title:    event.Title,
				Date:     event.Date,
				Duration: event.Duration,
				UserID:   event.UserID,
			})
		}
	}
	return results, nil
}

func (s *Storage) CheckEventOverlaps(ctx context.Context, date time.Time, duration string) (bool, error) {
	return false, nil

}
func New() (*Storage, error) {
	return &Storage{
		count:  1,
		events: make(EventsMap),
	}, nil
}
