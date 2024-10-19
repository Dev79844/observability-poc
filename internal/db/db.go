package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct{
	Pool *pgxpool.Pool
}

func InitDB() (*DB) {
	config, err := pgxpool.ParseConfig(os.Getenv("DB_URI"))
	if err!=nil{
		slog.Error("error parsing connection url", slog.Any("error", err))
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err!=nil{
		slog.Error("error creating a connection pool", slog.Any("error",err))
		return nil
	}

	fmt.Println("pool",pool.Ping(context.Background()))

	slog.Info("db connection established")
	return &DB{Pool: pool}
}

func (db *DB) Close(){
	db.Pool.Close()
}

func (db *DB) Stats() *pgxpool.Stat {
	return db.Pool.Stat()
}