package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/Dev79844/observeability-poc/internal/metrics"
	"github.com/Dev79844/observeability-poc/internal/middleware"
	"github.com/Dev79844/observeability-poc/internal/db"
)

func main(){
	database := db.InitDB()
	defer database.Close()
	metrics.InitMetrics(database.Pool)
	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.PrometheusMiddleware)

	router.Path("/metrics").Handler(promhttp.Handler())

	slog.Info("server started on port 9000")
	err := http.ListenAndServe(":9000", router)
	log.Fatal(err)
}