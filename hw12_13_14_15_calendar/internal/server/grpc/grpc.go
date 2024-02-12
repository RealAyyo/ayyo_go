package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/config"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/server"
	calendarV1 "github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/server/grpc/RealAyyo.hw12_13_14_15_calendar"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ServerAPI struct {
	calendarV1.UnimplementedCalendarServer
	logger     server.Logger
	app        Application
	GRPCServer *grpc.Server
	port       string
}

type Application interface {
	UpdateEvent(ctx context.Context, event *storage.Event) error
	DeleteEvent(ctx context.Context, eventID int, userID int) error
	CreateEvent(ctx context.Context, event *storage.Event) (*storage.Event, error)
	GetEventsByRange(ctx context.Context, userID int, dateFrom int64, dateTo int64) ([]storage.Event, error)
}

func NewServer(logger server.Logger, app Application, config config.GRPCConf) *ServerAPI {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor))
	serverAPI := &ServerAPI{
		logger:     logger,
		app:        app,
		GRPCServer: grpcServer,
		port:       config.Port,
	}

	calendarV1.RegisterCalendarServer(grpcServer, serverAPI)
	return serverAPI
}

func (s *ServerAPI) CreateEvent(ctx context.Context, req *calendarV1.CreateEventRequest) (*calendarV1.EventResponse, error) {
	event := &storage.Event{
		Title:    req.GetTitle(),
		Date:     req.GetDate().AsTime(),
		Duration: req.GetDuration(),
		UserID:   int(req.GetUserID()),
	}

	newEvent, err := s.app.CreateEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	return &calendarV1.EventResponse{
		Event: &calendarV1.Event{
			ID:       int32(newEvent.ID),
			Title:    newEvent.Title,
			Date:     timestamppb.New(newEvent.Date),
			Duration: newEvent.Duration,
			UserID:   int32(newEvent.UserID),
		},
	}, nil
}

func (s *ServerAPI) UpdateEvent(ctx context.Context, req *calendarV1.UpdateEventRequest) (*calendarV1.SuccessResponse, error) {
	event := &storage.Event{
		ID:       int(req.GetID()),
		Title:    req.GetTitle(),
		Date:     req.GetDate().AsTime(),
		Duration: req.GetDuration(),
		UserID:   int(req.GetUserID()),
	}

	err := s.app.UpdateEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	return &calendarV1.SuccessResponse{Message: "Event updated successfully"}, nil
}

func (s *ServerAPI) DeleteEvent(ctx context.Context, req *calendarV1.DeleteEventRequest) (*calendarV1.IdEventResponse, error) {
	err := s.app.DeleteEvent(ctx, int(req.GetID()), int(req.GetUserID()))
	if err != nil {
		return nil, err
	}

	return &calendarV1.IdEventResponse{ID: req.GetID()}, nil
}

func (s *ServerAPI) ListEvents(ctx context.Context, req *calendarV1.ListEventsRequest) (*calendarV1.ListEventsResponse, error) {
	events, err := s.app.GetEventsByRange(ctx, int(req.GetUserID()), req.GetDateFrom(), req.GetDateTo())
	if err != nil {
		return nil, err
	}
	eventsTransform := make([]*calendarV1.Event, 0, len(events))
	for _, event := range events {
		eventsTransform = append(eventsTransform, &calendarV1.Event{
			ID:       int32(event.ID),
			Title:    event.Title,
			Date:     timestamppb.New(event.Date),
			Duration: event.Duration,
			UserID:   int32(event.UserID),
		})
	}
	return &calendarV1.ListEventsResponse{Events: eventsTransform}, nil
}

func (s *ServerAPI) Stop() {
	s.GRPCServer.Stop()
}

func (s *ServerAPI) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := s.GRPCServer.Serve(lis); err != nil {
		return err
	}

	s.logger.Info(fmt.Sprintf("grpc server started. Address: %v", lis.Addr()))

	return nil
}
