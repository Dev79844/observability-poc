package middleware

import (
	"net/http"
	"strconv"
	"time"

    "github.com/Dev79844/observeability-poc/internal/metrics"
)

func PrometheusMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.TotalRequests.Inc()
		metrics.ActiveRequests.Inc()
		defer metrics.ActiveRequests.Dec()
        start := time.Now()

        rw := &responseWriter{w, http.StatusOK}
        next.ServeHTTP(rw, r)

        duration := time.Since(start).Seconds()
        metrics.HTTPRequestDuration.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(rw.statusCode)).Observe(duration)
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}