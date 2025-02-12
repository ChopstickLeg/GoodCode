package db

import (
	"context"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"
)

var (
	db     *pgx.Conn
	dbOnce sync.Once
)

func GetDB() (*pgx.Conn, error) {
	var err error
	dbOnce.Do(func() {
		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			err = errors.New("DATABASE_URL environment variable not set")

		}

		conn, err := pgx.Connect(context.Background(), dbURL)
		if err != nil {
			log.Fatal("Failed to connect to Neon Postgres:", err)
		}

		db = conn
	})

	return db, err
}
