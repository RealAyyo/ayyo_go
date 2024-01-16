package storage

import "time"

type Notification struct {
	EventId string
	Title   string
	Date    time.Time
	UserId  string
}
