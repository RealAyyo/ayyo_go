package internalhttp

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/config"
	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/controllers"
)

const (
	Timeout = 2
)

type Server struct {
	httpServer *http.Server
}

func NewServer(eventController *controllers.EventController, config config.HTTPConf) *Server {
	addr := net.JoinHostPort(config.Host, config.Port)
	httpServer := &http.Server{Addr: addr, ReadHeaderTimeout: Timeout * time.Second}

	http.Handle("/", loggingMiddleware(http.HandlerFunc(eventController.GetEventsByRange)))
	http.Handle("/create", loggingMiddleware(http.HandlerFunc(eventController.CreateEvent)))
	http.Handle("/update", loggingMiddleware(http.HandlerFunc(eventController.UpdateEvent)))
	http.Handle("/delete", loggingMiddleware(http.HandlerFunc(eventController.DeleteEvent)))

	return &Server{
		httpServer: httpServer,
	}
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.httpServer.Shutdown(ctx)
	return err
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
