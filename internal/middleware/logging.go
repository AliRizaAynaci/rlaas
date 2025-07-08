package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/AliRizaAynaci/rlaas/internal/logging"
)

// responseRecorder wraps http.ResponseWriter to record status and body size.
type responseRecorder struct {
	http.ResponseWriter     // embed the real ResponseWriter
	status              int // captures HTTP status code
	size                int // accumulates bytes written
}

// WriteHeader captures the status code before sending headers.
func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Write captures the number of bytes written to the client.
func (r *responseRecorder) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.size += n
	return n, err
}

// LoggingMiddleware logs HTTP requests with method, path, status, size, duration, and client IP.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// enrich base logger with HTTP component tag
		logger := logging.L.With(slog.String("component", "http"))
		if reqID := r.Header.Get("X-Request-ID"); reqID != "" {
			logger = logger.With(slog.String("request_id", reqID))
		}

		// wrap the ResponseWriter to capture status and size
		rec := &responseRecorder{
			ResponseWriter: w,
			status:         http.StatusOK, // default if WriteHeader never called
		}

		// call the next handler using our recorder
		next.ServeHTTP(rec, r)

		// after the handler finishes, log everything
		duration := time.Since(start)
		logger.Info("HTTP request completed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rec.status),
			slog.Int("size", rec.size),
			slog.Duration("duration", duration),
			slog.String("remote_addr", r.RemoteAddr),
		)
	})
}
