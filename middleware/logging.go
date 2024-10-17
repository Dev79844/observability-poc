package middleware

import(
	"net/http"
	"time"
	"log/slog"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w,r)
		slog.Info("request recorded",slog.Any("method", r.Method), slog.Any("path", r.URL.Path), slog.Any("duration", time.Since(start)))
	})
} 