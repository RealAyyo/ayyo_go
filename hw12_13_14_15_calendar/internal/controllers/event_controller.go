package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage"
)

type EventController struct {
	app       Application
	validator Validator
}

type Validator interface {
	Validate(method string, r *http.Request, data interface{}) error
}

type Application interface {
	UpdateEvent(ctx context.Context, event *storage.Event) (int, error)
	DeleteEvent(ctx context.Context, eventID int, userID int) error
	CreateEvent(ctx context.Context, event *storage.Event) (int, error)
	GetEventsByRange(ctx context.Context, userID int, dateFrom int64, dateTo int64) ([]storage.Event, error)
}

func NewEventController(app Application, validator Validator) *EventController {
	return &EventController{
		app:       app,
		validator: validator,
	}
}

func (e *EventController) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var dto CreateEventDto
	err := e.validator.Validate("POST", r, &dto)
	if err != nil {
		resp := ErrorResponse{
			Message: err.Error(),
			Err:     ErrHas,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	date := time.Unix(dto.Date, 0)

	event := &storage.Event{
		Title:    dto.Title,
		Date:     date,
		Duration: dto.Duration,
		UserID:   dto.UserID,
	}

	id, err := e.app.CreateEvent(ctx, event)
	if err != nil {
		resp := ErrorResponse{
			Message: err.Error(),
			Err:     ErrHas,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := SuccessResponse{
		Message: "Event created successfully",
		Data: &IdResponseDto{
			ID: id,
		},
		Err: ErrNo,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (e *EventController) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var dto UpdateEventDto
	err := e.validator.Validate("PATCH", r, &dto)
	if err != nil {
		resp := ErrorResponse{
			Message: err.Error(),
			Err:     ErrHas,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	date := time.Unix(dto.Date, 0)

	event := &storage.Event{
		Title:    dto.Title,
		Date:     date,
		Duration: dto.Duration,
		UserID:   dto.UserID,
		ID:       dto.ID,
	}

	id, err := e.app.UpdateEvent(ctx, event)
	if err != nil {
		resp := ErrorResponse{
			Message: err.Error(),
			Err:     ErrHas,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := SuccessResponse{
		Message: "Event updated successfully",
		Data: &IdResponseDto{
			ID: id,
		},
		Err: ErrNo,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (e *EventController) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var dto DeleteEventDto
	err := e.validator.Validate("POST", r, &dto)
	if err != nil {
		resp := ErrorResponse{
			Message: err.Error(),
			Err:     ErrHas,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	err = e.app.DeleteEvent(ctx, dto.ID, dto.UserID)
	if err != nil {
		resp := ErrorResponse{
			Message: err.Error(),
			Err:     ErrHas,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
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
	json.NewEncoder(w).Encode(resp)
}

func (e *EventController) GetEventsByRange(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var dto GetEventsDto
	err := e.validator.Validate("GET", r, &dto)
	if err != nil {
		resp := ErrorResponse{
			Message: err.Error(),
			Err:     ErrHas,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	parsedUserID, err := strconv.ParseInt(dto.UserID, 10, 64)
	parsedDataFrom, err := strconv.ParseInt(dto.DateFrom, 10, 64)
	parsedDataTo, err := strconv.ParseInt(dto.DateTo, 10, 64)
	if err != nil {
		resp := ErrorResponse{
			Message: err.Error(),
			Err:     ErrHas,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	dateFrom := time.Unix(parsedDataFrom, 0)
	dateTo := time.Unix(parsedDataTo, 0)

	fmt.Println(dateFrom)
	fmt.Println(dateTo)

	events, err := e.app.GetEventsByRange(ctx, int(parsedUserID), parsedDataFrom, parsedDataTo)
	if err != nil {
		resp := ErrorResponse{
			Message: err.Error(),
			Err:     ErrHas,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := SuccessResponse{
		Message: "Event received successfully",
		Data:    events,
		Err:     ErrNo,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
