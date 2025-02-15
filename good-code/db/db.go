package db

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var (
	db     *sql.DB
	dbOnce sync.Once
)

func GetDB() (*sql.DB, error) {
	var err error
	dbOnce.Do(func() {
		dbURL := os.Getenv("DATABASE_DATABASE_URL")
		if dbURL == "" {
			err = errors.New("DATABASE_URL environment variable not set")
		}

		conn, err := sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatal("Failed to connect to Neon Postgres:", err)
		}

		db = conn
	})

	return db, err
}
