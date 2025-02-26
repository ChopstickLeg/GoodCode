package db

import (
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

func GetDB() (*gorm.DB, error) {
	var err error

	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
			Colorful:                  true,
		},
	)

	dbOnce.Do(func() {
		dbURL := os.Getenv("DATABASE_DATABASE_URL")
		if dbURL == "" {
			err = errors.New("DATABASE_URL environment variable not set")
		}

		conn, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
			Logger: logger,
		})
		if err != nil {
			log.Fatal("Failed to connect to Neon Postgres:", err)
		}

		db = conn
	})

	return db, err
}
