package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/app"
	storage2 "github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

type TestCases struct {
	name        string
	event       *storage2.Event
	errRequired error
}

var testCases = []TestCases{
	{
		name: "First Event",
		event: &storage2.Event{
			UserId:   2,
			Title:    "Meet",
			Duration: "1:00:00",
			Date:     time.Now(),
		},
		errRequired: nil,
	},
	{
		name: "Second Event",
		event: &storage2.Event{
			Title:    "Daily",
			Duration: "0:30:00",
			Date:     time.Now(),
		},
		errRequired: app.ErrUserIdRequired,
	},
	{
		name: "Third Event",
		event: &storage2.Event{
			UserId:   2,
			Duration: "2:00:00",
			Date:     time.Now(),
		},
		errRequired: app.ErrTitleRequired,
	},
	{
		name: "Fourth Event",
		event: &storage2.Event{
			UserId: 3,
			Title:  "Daily",
			Date:   time.Now(),
		},
		errRequired: app.ErrDurationRequired,
	},
	{
		name: "Fourth Event",
		event: &storage2.Event{
			UserId:   3,
			Title:    "Daily",
			Duration: "2:00:00",
		},
		errRequired: app.ErrDateRequired,
	},
}

func TestStorage(t *testing.T) {
	ctx := context.Background()

	storage, err := New()
	require.NoError(t, err)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := storage.AddEvent(ctx, testCase.event)
			if testCase.event != nil {
				require.ErrorIs(t, err, testCase.errRequired)
			}
		})
	}

	t.Run("List Events", func(t *testing.T) {
		listEvents, err := storage.ListEvents(ctx, 2, time.Now().Add(-time.Minute), time.Now().AddDate(0, 0, 1))
		if err != nil {
			fmt.Println(err)
		}

		require.Equal(t, len(listEvents), 1)
	})

	t.Run("Update Event", func(t *testing.T) {
		newTitle := "New Title 12345678910"
		updated := &storage2.Event{
			ID:     1,
			UserId: 2,
			Title:  newTitle,
		}
		err := storage.UpdateEvent(ctx, updated)

		listEvents, err := storage.ListEvents(ctx, 2, time.Now().Add(-time.Minute), time.Now().AddDate(0, 0, 1))
		if err != nil {
			fmt.Println(err)
		}

		require.Equal(t, listEvents[0].Title, newTitle)
	})

	t.Run("Delete Event", func(t *testing.T) {
		err := storage.DeleteEvent(ctx, 1, 2)

		listEvents, err := storage.ListEvents(ctx, 2, time.Now().Add(-time.Minute), time.Now().AddDate(0, 0, 1))
		if err != nil {
			fmt.Println(err)
		}
		require.Equal(t, len(listEvents), 0)
	})

	t.Run("Concurrent Adding", func(t *testing.T) {
		var wg sync.WaitGroup
		numberOfGoroutines := 35

		for i := 1; i < numberOfGoroutines+1; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				event := &storage2.Event{
					UserId:   2,
					Title:    fmt.Sprintf("Event %d", i),
					Duration: "1:00:00",
					Date:     time.Now(),
				}
				err := storage.AddEvent(ctx, event)
				require.NoError(t, err)
			}(i)
		}

		wg.Wait()

		listEvents, err := storage.ListEvents(ctx, 2, time.Now().Add(-time.Hour), time.Now().Add(time.Hour))
		require.NoError(t, err)
		require.Equal(t, numberOfGoroutines, len(listEvents))
	})
}
