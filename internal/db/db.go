package db

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct{
	Pool *pgxpool.Pool
}

func InitDB() (*DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(os.Getenv("DB_URI"))
	if err!=nil{
		slog.Error("error parsing connection url", slog.Any("error", err))
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err!=nil{
		slog.Error("error creating a connection pool", slog.Any("error",err))
		return nil
	}

	if err := pool.Ping(ctx); err != nil {
		slog.Error("error connecting to database", slog.Any("error", err))
		pool.Close()
		return nil
	}

	db := &DB{Pool: pool}

	if err := db.migrate(ctx); err != nil {
		slog.Error("migration failed", slog.Any("error", err))
		pool.Close()
		return nil
	}

	slog.Info("db connection established")
	
	return db
}

func (db *DB)migrate(ctx context.Context) error {
	tx, err := db.Pool.Begin(ctx)
	if err!=nil{
		return err
	}

	defer tx.Rollback(ctx)


	_, err = tx.Exec(ctx, 
					`CREATE TABLE IF NOT EXISTS todos (
						id text primary key,
						title text,
						created_at timestamp);`,
					)
	if err!=nil{
		slog.Error("error creating todos table", slog.Any("error", err))
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error("error committing migration transaction", slog.Any("error", err))
		return err
	}

	slog.Info("tables migrated")
	return nil
}

func (db *DB) Close(){
	db.Pool.Close()
}

func (db *DB) Stats() *pgxpool.Stat {
	return db.Pool.Stat()
}