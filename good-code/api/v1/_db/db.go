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
		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			err = errors.New("DATABASE_URL environment variable not set")
			return
		}

		conn, dbErr := gorm.Open(postgres.Open(dbURL), &gorm.Config{
			Logger: logger,
		})
		if dbErr != nil {
			err = dbErr
			log.Printf("Failed to connect to Neon Postgres: %v", dbErr)
			return
		}

		db = conn
	})

	if err != nil {
		return nil, err
	}

	if db == nil {
		return nil, errors.New("database connection not initialized")
	}

	return db, nil
}
