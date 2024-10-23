package api

import (
	"github.com/Dev79844/observeability-poc/internal/db"
	"github.com/Dev79844/observeability-poc/internal/middleware"
	"github.com/gorilla/mux"
)

func SetApiRoutes(r *mux.Router, db *db.DB) {
	h := &Handler{DB: db}

	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.PrometheusMiddleware)

	r.HandleFunc("/todos", h.CreateTodo).Methods("POST")
	r.HandleFunc("/todo/{id}", h.GetTodo).Methods("GET")
	r.HandleFunc("/todo/{id}", h.UpdateTodo).Methods("PUT")
	r.HandleFunc("/todo/{id}", h.DeleteTodo).Methods("DELETE")
}