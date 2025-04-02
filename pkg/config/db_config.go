package config

import (
	"context"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var (
	dbInstance *pgxpool.Pool
	dbOnce     sync.Once
)

func GetDBConnection() *pgxpool.Pool {
	dbOnce.Do(func() {
		dbURL := GetEnv("DATABASE_URL", "")
		if dbURL == "" {
			log.Fatal().Msg("DATABASE_URL environment variable is not set")
		}

		poolConfig, err := pgxpool.ParseConfig(dbURL)
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to parse database URL")
		}

		dbInstance, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to connect to database")
		}

		if err := dbInstance.Ping(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("Unable to ping database")
		}

		log.Info().Msg("Successfully connected to database")
	})

	return dbInstance
}

func CloseDBConnection() {
	if dbInstance != nil {
		dbInstance.Close()
	}
}

func GetEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
