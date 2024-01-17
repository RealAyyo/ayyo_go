package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/config"
	storage2 "github.com/RealAyyo/ayyo_go/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	httpServer *http.Server
	logger     Logger
	app        Application
}

type Logger interface {
	Info(msg string, attrs ...any)
	Error(msg string, attrs ...any)
	Debug(msg string, attrs ...any)
	Warn(msg string, attrs ...any)
}

type Application interface {
	CreateEvent(ctx context.Context, event *storage2.Event) error
	GetEventsForRange(ctx context.Context, userId int, dateFrom time.Time, dateRange int) ([]storage2.Event, error)
}

func NewServer(logger Logger, app Application, config config.HttpConf) *Server {
	addr := net.JoinHostPort(config.Host, config.Port)
	httpServer := &http.Server{Addr: addr}

	http.HandleFunc("/", helloHandler)

	return &Server{
		logger:     logger,
		app:        app,
		httpServer: httpServer,
	}
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.httpServer.Shutdown(ctx)
	return err
}
func helloHandler(w http.ResponseWriter, r *http.Request) {
	ip := ReadUserIP(r)
	fmt.Println(fmt.Sprintf(
		"%v [%v] %v %v %v %v %v %v",
		ip,
		time.Now().Format("02/Jan/2006:15:04:05 -0700"),
		r.Method,
		r.URL.Path,
		r.Proto,
		200,
		100,
		r.UserAgent(),
	))
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
