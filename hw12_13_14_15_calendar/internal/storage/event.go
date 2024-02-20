package storage

import (
	"time"
)

type Event struct {
	ID               int
	Title            string
	Date             time.Time
	Duration         string
	UserID           int
	Description      string
	NotificationTime string
}
