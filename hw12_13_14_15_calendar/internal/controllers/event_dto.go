package controllers

import "github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage"

type CreateEventDto struct {
	Title            string `json:"title" validate:"required"`
	Date             int64  `json:"date" validate:"required"`
	Duration         string `json:"duration" validate:"required"`
	UserID           int    `json:"user_id" validate:"required"`
	Description      string `json:"description"`
	NotificationTime string `json:"notification_time"`
}

type UpdateEventDto struct {
	Title            string `json:"title"`
	Date             int64  `json:"date"`
	Duration         string `json:"duration"`
	Description      string `json:"description"`
	NotificationTime string `json:"notification_time"`
	ID               int    `json:"id" validate:"required"`
	UserID           int    `json:"user_id" validate:"required"`
}

type DeleteEventDto struct {
	ID     int `json:"id" validate:"required"`
	UserID int `json:"user_id" validate:"required"`
}

type GetEventsDto struct {
	UserID   string `json:"user_id" validate:"required"`
	DateFrom string `json:"date_from" validate:"required"`
	DateTo   string `json:"date_to" validate:"required"`
}

type IdResponseDto struct {
	ID int `json:"id" validate:"required"`
}
type EventResponseDto struct {
	Event storage.Event `json:"event"`
}