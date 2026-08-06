package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"

	"myapp/internal/logger"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := newRequestID()

		ctx := logger.WithRequestID(r.Context(), requestID)
		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", requestID)

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		attrs := []any{
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"duration_ms", duration.Milliseconds(),
			"request_id", requestID,
			"remote_addr", r.RemoteAddr,
		}

		if userID, ok := r.Context().Value(UserIDContextKey{}).(int); ok && userID != 0 {
			attrs = append(attrs, "user_id", userID)
		}

		level := slog.LevelInfo
		if rw.status >= 500 {
			level = slog.LevelError
		} else if rw.status >= 400 {
			level = slog.LevelWarn
		}

		slog.Log(ctx, level, "http_request", attrs...)
	})
}

func newRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return hex.EncodeToString([]byte(time.Now().Format("150405.000")))
	}
	return hex.EncodeToString(b)
}
