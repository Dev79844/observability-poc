package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Dev79844/observeability-poc/internal/db"
	"github.com/gorilla/mux"
)

type Handler struct {
	DB *db.DB
}

func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request){
	var todo struct{
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		slog.Error("error reading the body", slog.Any("error", err))
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	resp, err := h.DB.CreateTodo(context.Background(), todo.Title)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetTodo(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    item, err := h.DB.GetTodo(context.Background(), vars["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(item)
}

func (h *Handler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    err := h.DB.DeleteTodo(context.Background(), vars["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode("deleted successfully")
}

func (h *Handler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
    var todo struct{
		Id    string `json:"id"`
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		slog.Error("error reading the body", slog.Any("error", err))
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    item, err := h.DB.UpdateTodo(context.Background(), todo.Id, todo.Title)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(item)
}

