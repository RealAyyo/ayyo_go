package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/validators"
)

type EventController struct {
	app    Application
	logger Logger
}

type Logger interface {
	Info(msg string, attrs ...any)
	Error(msg string, attrs ...any)
	Debug(msg string, attrs ...any)
	Warn(msg string, attrs ...any)
}

type Application interface {
	UpdateEvent(ctx context.Context, event *storage.Event) error
	DeleteEvent(ctx context.Context, eventID int, userID int) error
	CreateEvent(ctx context.Context, event *storage.Event) (*storage.Event, error)
	GetEventsByRange(ctx context.Context, userID int, dateFrom int64, dateTo int64) ([]storage.Event, error)
}

func NewEventController(app Application, logger Logger) *EventController {
	return &EventController{
		app:    app,
		logger: logger,
	}
}

func (e *EventController) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var dto CreateEventDto
	err := validators.Validate("POST", r, &dto)
	if err != nil {
		e.logger.Error(err.Error())
		sendErrorResponse(err, w)
		return
	}

	date := time.Unix(dto.Date, 0)

	event := &storage.Event{
		Title:            dto.Title,
		Date:             date,
		Duration:         dto.Duration,
		UserID:           dto.UserID,
		Description:      dto.Description,
		NotificationTime: dto.NotificationTime,
	}

	newEvent, err := e.app.CreateEvent(ctx, event)
	if err != nil {
		e.logger.Error(err.Error())
		sendErrorResponse(err, w)
		return
	}

	resp := SuccessResponse{
		Message: "Event created successfully",
		Data: &EventResponseDto{
			Event: *newEvent,
		},
		Err: ErrNo,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		e.logger.Error(ErrEncodeJson.Error(), ErrEncodeJson)
		return
	}
}

func (e *EventController) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var dto UpdateEventDto
	err := validators.Validate("PATCH", r, &dto)
	if err != nil {
		e.logger.Error(err.Error())
		sendErrorResponse(err, w)
		return
	}

	date := time.Unix(dto.Date, 0)

	event := &storage.Event{
		Title:            dto.Title,
		Date:             date,
		Duration:         dto.Duration,
		UserID:           dto.UserID,
		ID:               dto.ID,
		Description:      dto.Description,
		NotificationTime: dto.NotificationTime,
	}

	err = e.app.UpdateEvent(ctx, event)
	if err != nil {
		e.logger.Error(err.Error())
		sendErrorResponse(err, w)
		return
	}

	resp := SuccessResponse{
		Message: "Event updated successfully",
		Err:     ErrNo,
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		e.logger.Error(ErrEncodeJson.Error(), ErrEncodeJson)
		return
	}
}

func (e *EventController) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var dto DeleteEventDto
	err := validators.Validate("POST", r, &dto)
	if err != nil {
		e.logger.Error(err.Error())
		sendErrorResponse(err, w)
		return
	}

	err = e.app.DeleteEvent(ctx, dto.ID, dto.UserID)
	if err != nil {
		e.logger.Error(err.Error())
		sendErrorResponse(err, w)
		return
	}

	resp := SuccessResponse{
		Message: "Event deleted successfully",
		Data: &IdResponseDto{
			ID: dto.ID,
		},
		Err: ErrNo,
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		e.logger.Error(ErrEncodeJson.Error(), ErrEncodeJson)
		return
	}
}

func (e *EventController) GetEventsByRange(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var dto GetEventsDto
	err := validators.Validate("GET", r, &dto)
	if err != nil {
		e.logger.Error(err.Error())
		sendErrorResponse(err, w)
		return
	}

	parsedUserID, err := strconv.ParseInt(dto.UserID, 10, 64)
	if err != nil {
		e.logger.Error(err.Error())
		sendErrorResponse(err, w)
		return
	}

	parsedDataFrom, err := strconv.ParseInt(dto.DateFrom, 10, 64)
	if err != nil {
		e.logger.Error(err.Error())
		sendErrorResponse(err, w)
		return
	}

	parsedDataTo, err := strconv.ParseInt(dto.DateTo, 10, 64)
	if err != nil {
		e.logger.Error(err.Error())
		sendErrorResponse(err, w)
		return
	}

	events, err := e.app.GetEventsByRange(ctx, int(parsedUserID), parsedDataFrom, parsedDataTo)
	if err != nil {
		e.logger.Error(err.Error())
		sendErrorResponse(err, w)
		return
	}

	resp := SuccessResponse{
		Message: "Event received successfully",
		Data:    events,
		Err:     ErrNo,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		e.logger.Error(ErrEncodeJson.Error(), ErrEncodeJson)
		return
	}
}
