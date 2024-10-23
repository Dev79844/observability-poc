package db

import (
	"context"
	"log/slog"
	"time"

	"github.com/Dev79844/observeability-poc/internal/metrics"
	"github.com/google/uuid"
)
	

type Todo struct{
	Id			string		`json:"id"`
	Title 		string		`json:"title"`
	CreatedAt	time.Time	`json:"created_at"`
}

func (db *DB) CreateTodo(ctx context.Context, title string) (*Todo, error) {
	timer := metrics.NewDBTimer("create_item")
	defer timer.ObserveDuration()
	id := uuid.NewString()
	created_at := time.Now()
	var todo Todo
	err := db.Pool.QueryRow(ctx, 
					"INSERT INTO todos(id, title, created_at) VALUES($1, $2, $3) RETURNING id, title, created_at", 
					id, title, created_at).Scan(&todo.Id, &todo.Title, &todo.CreatedAt)

	if err!=nil{
		metrics.DBErrors.WithLabelValues("create_item").Inc()
		slog.Error("error executing insert query", slog.Any("error", err))
		return nil, err
	}

	return &todo, nil
}

func (db *DB) GetTodo(ctx context.Context, id string) (*Todo, error) {
    timer := metrics.NewDBTimer("get_item")
    defer timer.ObserveDuration()

    var todo Todo
    err := db.Pool.QueryRow(ctx,
        "SELECT id, title, created_at FROM todos WHERE id = $1",
        id,
    ).Scan(&todo.Id, &todo.Title, &todo.CreatedAt)

    if err != nil {
        metrics.DBErrors.WithLabelValues("get_item").Inc()
        return nil, err
    }
    return &todo, nil
}

func (db *DB) UpdateTodo(ctx context.Context, id string, name string) (*Todo, error) {
    timer := metrics.NewDBTimer("update_item")
    defer timer.ObserveDuration()

    var todo Todo
    err := db.Pool.QueryRow(ctx,
        "UPDATE todos SET title = $1 WHERE id = $2 RETURNING id, title, created_at",
        name, id,
    ).Scan(&todo.Id, &todo.Title, &todo.CreatedAt)

    if err != nil {
        metrics.DBErrors.WithLabelValues("update_item").Inc()
        return nil, err
    }
    return &todo, nil
}

func (db *DB) DeleteTodo(ctx context.Context, id string) error {
    timer := metrics.NewDBTimer("delete_item")
    defer timer.ObserveDuration()

    commandTag, err := db.Pool.Exec(ctx,
        "DELETE FROM todos WHERE id = $1",
        id,
    )
    if err != nil {
        metrics.DBErrors.WithLabelValues("delete_item").Inc()
        return err
    }

    if commandTag.RowsAffected() == 0 {
        metrics.DBErrors.WithLabelValues("delete_item_not_found").Inc()
        return err
    }
    return nil
}