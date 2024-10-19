package metrics

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
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

type DBStatsCollector struct {
	acquireCount     	*prometheus.Desc
    acquireDuration  	*prometheus.Desc
    acquiredConns    	*prometheus.Desc
    canceledAcquires 	*prometheus.Desc
    constructingConns 	*prometheus.Desc
    idleConns        	*prometheus.Desc
    maxConns         	*prometheus.Desc
    totalConns       	*prometheus.Desc
    pool             	interface{ Stat() *pgxpool.Stat }
}

func NewDBStatsCollector(pool interface{ Stat() *pgxpool.Stat }) *DBStatsCollector {
    return &DBStatsCollector{
        acquireCount: prometheus.NewDesc(
            "pgx_pool_acquire_count",
            "Total number of connection acquires",
            nil, nil,
        ),
        acquireDuration: prometheus.NewDesc(
            "pgx_pool_acquire_duration_seconds",
            "Total duration of connection acquires",
            nil, nil,
        ),
        acquiredConns: prometheus.NewDesc(
            "pgx_pool_acquired_conns",
            "Number of currently acquired connections",
            nil, nil,
        ),
        canceledAcquires: prometheus.NewDesc(
            "pgx_pool_canceled_acquires",
            "Number of canceled acquires",
            nil, nil,
        ),
        constructingConns: prometheus.NewDesc(
            "pgx_pool_constructing_conns",
            "Number of connections being constructed",
            nil, nil,
        ),
        idleConns: prometheus.NewDesc(
            "pgx_pool_idle_conns",
            "Number of idle connections",
            nil, nil,
        ),
        maxConns: prometheus.NewDesc(
            "pgx_pool_max_conns",
            "Maximum number of connections",
            nil, nil,
        ),
        totalConns: prometheus.NewDesc(
            "pgx_pool_total_conns",
            "Total number of connections",
            nil, nil,
        ),
        pool: pool,
    }
}

func (c *DBStatsCollector) Describe(ch chan<- *prometheus.Desc) {
    ch <- c.acquireCount
    ch <- c.acquireDuration
    ch <- c.acquiredConns
    ch <- c.canceledAcquires
    ch <- c.constructingConns
    ch <- c.idleConns
    ch <- c.maxConns
    ch <- c.totalConns
}

func (c *DBStatsCollector) Collect(ch chan<- prometheus.Metric) {
    stats := c.pool.Stat()
    ch <- prometheus.MustNewConstMetric(c.acquireCount, prometheus.CounterValue, float64(stats.AcquireCount()))
    ch <- prometheus.MustNewConstMetric(c.acquireDuration, prometheus.CounterValue, stats.AcquireDuration().Seconds())
    ch <- prometheus.MustNewConstMetric(c.acquiredConns, prometheus.GaugeValue, float64(stats.AcquiredConns()))
    ch <- prometheus.MustNewConstMetric(c.canceledAcquires, prometheus.CounterValue, float64(stats.CanceledAcquireCount()))
    ch <- prometheus.MustNewConstMetric(c.constructingConns, prometheus.GaugeValue, float64(stats.ConstructingConns()))
    ch <- prometheus.MustNewConstMetric(c.idleConns, prometheus.GaugeValue, float64(stats.IdleConns()))
    ch <- prometheus.MustNewConstMetric(c.maxConns, prometheus.GaugeValue, float64(stats.MaxConns()))
    ch <- prometheus.MustNewConstMetric(c.totalConns, prometheus.GaugeValue, float64(stats.TotalConns()))
}

func InitMetrics(pool *pgxpool.Pool) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		HTTPRequestDuration,
		TotalRequests,
		ActiveRequests,
		ResponseStatus,
		NewDBStatsCollector(pool),
	)
	slog.Info("metrics initialised")
}
