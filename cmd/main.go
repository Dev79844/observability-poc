package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/Dev79844/observeability-poc/internal/metrics"
	"github.com/Dev79844/observeability-poc/internal/middleware"
)

type responseWriter struct{
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int){
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func main(){
	metrics.InitMetrics()
	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.PrometheusMiddleware)

	router.Path("/metrics").Handler(promhttp.Handler())

	slog.Info("server started on port 9000")
	err := http.ListenAndServe(":9000", router)
	log.Fatal(err)
}