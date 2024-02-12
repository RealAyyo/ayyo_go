package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/app"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

type TestCases struct {
	name        string
	event       *storage.Event
	errRequired error
}

var testCases = []TestCases{
	{
		name: "First Event",
		event: &storage.Event{
			UserID:   2,
			Title:    "Meet",
			Duration: "1:00:00",
			Date:     time.Now(),
		},
		errRequired: nil,
	},
	{
		name: "Second Event",
		event: &storage.Event{
			Title:    "Daily",
			Duration: "0:30:00",
			Date:     time.Now(),
		},
		errRequired: app.ErrUserIDRequired,
	},
	{
		name: "Third Event",
		event: &storage.Event{
			UserID:   2,
			Duration: "2:00:00",
			Date:     time.Now(),
		},
		errRequired: app.ErrTitleRequired,
	},
	{
		name: "Fourth Event",
		event: &storage.Event{
			UserID: 3,
			Title:  "Daily",
			Date:   time.Now(),
		},
		errRequired: app.ErrDurationRequired,
	},
	{
		name: "Fourth Event",
		event: &storage.Event{
			UserID:   3,
			Title:    "Daily",
			Duration: "2:00:00",
		},
		errRequired: app.ErrDateRequired,
	},
}

func TestStorage(t *testing.T) {
	ctx := context.Background()

	storageService, err := New()
	require.NoError(t, err)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := storageService.AddEvent(ctx, testCase.event)
			if testCase.errRequired != nil {
				require.ErrorIs(t, err, testCase.errRequired)
			}
		})
	}

	t.Run("List Events", func(t *testing.T) {
		listEvents, err := storageService.ListEvents(ctx, 2, time.Now().Add(-time.Minute), time.Now().AddDate(0, 0, 1))
		require.NoError(t, err)

		require.Equal(t, len(listEvents), 1)
	})

	t.Run("Update Event", func(t *testing.T) {
		newTitle := "New Title 12345678910"
		updated := &storage.Event{
			ID:     1,
			UserID: 2,
			Title:  newTitle,
		}
		err := storageService.UpdateEvent(ctx, updated)
		require.NoError(t, err)

		listEvents, err := storageService.ListEvents(ctx, 2, time.Now().Add(-time.Minute), time.Now().AddDate(0, 0, 1))
		require.NoError(t, err)

		require.Equal(t, listEvents[0].Title, newTitle)
	})

	t.Run("Delete Event", func(t *testing.T) {
		err := storageService.DeleteEvent(ctx, 1, 2)
		require.NoError(t, err)

		listEvents, err := storageService.ListEvents(ctx, 2, time.Now().Add(-time.Minute), time.Now().AddDate(0, 0, 1))

		require.NoError(t, err)
		require.Equal(t, len(listEvents), 0)
	})

	t.Run("Concurrent Adding", func(t *testing.T) {
		var wg sync.WaitGroup
		numberOfGoroutines := 35

		for i := 1; i < numberOfGoroutines+1; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				event := &storage.Event{
					UserID:   2,
					Title:    fmt.Sprintf("Event %d", i),
					Duration: "1:00:00",
					Date:     time.Now(),
				}

				_, err := storageService.AddEvent(ctx, event)
				require.NoError(t, err)
			}(i)
		}

		wg.Wait()

		listEvents, err := storageService.ListEvents(ctx, 2, time.Now().Add(-time.Hour), time.Now().Add(time.Hour))
		require.NoError(t, err)
		require.Equal(t, numberOfGoroutines, len(listEvents))
	})
}
