package db

import (
	"errors"
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

func GetDB() (*gorm.DB, error) {
	var err error
	dbOnce.Do(func() {
		dbURL := os.Getenv("DATABASE_DATABASE_URL")
		if dbURL == "" {
			err = errors.New("DATABASE_URL environment variable not set")
		}

		conn, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to Neon Postgres:", err)
		}

		db = conn
	})

	return db, err
}
