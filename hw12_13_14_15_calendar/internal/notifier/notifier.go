package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage"
)

const (
	eventQueue = "calendar_event"
)

type App interface {
	GetEventsToNotify(ctx context.Context) ([]storage.Event, error)
}

type Broker interface {
	Send(ctx context.Context, queue string, message []byte) error
}

type Logger interface {
	Info(msg string, attrs ...any)
	Error(msg string, attrs ...any)
	Debug(msg string, attrs ...any)
	Warn(msg string, attrs ...any)
}

type Notification struct {
	EventID string
	Title   string
	Date    time.Time
	UserID  string
}

type Notifier struct {
	amqp           Broker
	app            App
	log            Logger
	intervalNotify time.Duration
}

func New(amqp Broker, app App, log Logger, intervalNotify string) (*Notifier, error) {
	interval, err := time.ParseDuration(intervalNotify)
	if err != nil {
		return nil, err
	}

	return &Notifier{
		amqp:           amqp,
		app:            app,
		log:            log,
		intervalNotify: interval,
	}, nil
}

func (n *Notifier) Start() {
	ticker := time.NewTicker(n.intervalNotify)
	defer ticker.Stop()
	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		events, err := n.app.GetEventsToNotify(ctx)
		if err != nil {
			n.log.Error("error getting events: " + err.Error())
			cancel()
			continue
		}

		for _, event := range events {
			var notification Notification
			notification.EventID = strconv.Itoa(event.ID)
			notification.Title = event.Title
			notification.Date = event.Date
			notification.UserID = strconv.Itoa(event.UserID)

			notificationBytes, err := json.Marshal(notification)
			fmt.Println(string(notificationBytes))
			if err != nil {
				n.log.Error("error marshaling notification: " + err.Error())
				continue
			}

			err = n.amqp.Send(context.Background(), eventQueue, notificationBytes)
			if err != nil {
				n.log.Error("error sending message: " + err.Error())
				continue
			}
		}
		cancel()
	}
}
