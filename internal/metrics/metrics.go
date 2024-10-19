package metrics

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var(
	TotalRequests = promauto.NewCounter(prometheus.CounterOpts{
        Name: "http_requests_total",
        Help: "Total number of HTTP requests.",
    })

	ActiveRequests = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "http_requests_active",
        Help: "Number of active HTTP requests.",
    })

	ResponseStatus = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_status",
			Help: "Status of HTTP response",
		},
		[]string{"status"},
	)
	
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_response_time_seconds",
			Help: "Duration of HTTP requests.",
		},
		[]string{"path", "method", "status_code"},
	)
)

func InitMetrics() {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		HTTPRequestDuration,
		TotalRequests,
		ActiveRequests,
		ResponseStatus,
	)
	slog.Info("metrics initialised")
}
