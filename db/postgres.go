package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect(connStr string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal("Unable to ping database:", err)
	}

	DB = pool
}
