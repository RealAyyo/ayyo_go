package internalhttp

import (
	"log"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		startTime := time.Now()
		next.ServeHTTP(rec, r)
		duration := time.Since(startTime)
		ip := ReadUserIP(r)
		log.Printf(
			"%v [%v] %v %v %v %v %v %v",
			ip,
			time.Now().Format(time.RFC822),
			r.Method,
			r.URL.Path,
			r.Proto,
			rec.status,
			duration,
			r.UserAgent(),
		)
	})
}
