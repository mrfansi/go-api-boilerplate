package middleware

import (
	"bufio"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type LoggerMiddleware struct{}

func NewLoggerMiddleware() *LoggerMiddleware {
	return &LoggerMiddleware{}
}

// Logger logs HTTP request details
func (m *LoggerMiddleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture the status code
		ww := &responseWriter{w: w, status: http.StatusOK}

		// Process request
		next.ServeHTTP(ww, r)

		// Calculate duration
		duration := time.Since(start)

		// Log request details
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Int("status", ww.status).
			Str("duration", duration.String()).
			Str("user_agent", r.UserAgent()).
			Msg("HTTP Request")
	})
}

// responseWriter is a custom response writer that captures the status code
type responseWriter struct {
	w      http.ResponseWriter
	status int
}

func (rw *responseWriter) Header() http.Header {
	return rw.w.Header()
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.w.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.w.WriteHeader(statusCode)
}

// Implement http.Pusher if the underlying ResponseWriter implements it
func (rw *responseWriter) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := rw.w.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return http.ErrNotSupported
}

// Implement http.Flusher if the underlying ResponseWriter implements it
func (rw *responseWriter) Flush() {
	if flusher, ok := rw.w.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Implement http.Hijacker if the underlying ResponseWriter implements it
func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := rw.w.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}
