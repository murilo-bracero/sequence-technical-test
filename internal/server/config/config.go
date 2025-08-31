package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/murilo-bracero/sequence-technical-test/internal/utils"
)

type Config struct {
	AppPort          string
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	MaxDbConnections int
	MinDbConnections int
	MaxConnIdleTime  int

	MaxSequencePagination int

	MaxCacheMemory  int
	CacheLifeWindow int
}

func New() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Warn("failed to load .env file", err.Error(), err)
	}

	return &Config{
		AppPort:          os.Getenv("APP_PORT"),
		PostgresHost:     os.Getenv("DB_HOST"),
		PostgresPort:     utils.SafeAtoi(os.Getenv("DB_PORT"), 5432),
		PostgresUser:     os.Getenv("DB_USER"),
		PostgresPassword: os.Getenv("DB_PASSWORD"),
		PostgresDatabase: os.Getenv("DB_NAME"),
		MaxDbConnections: utils.SafeAtoi(os.Getenv("DB_MAX_CONNECTIONS"), 10),
		MinDbConnections: utils.SafeAtoi(os.Getenv("DB_MIN_CONNECTIONS"), 1),
		MaxConnIdleTime:  utils.SafeAtoi(os.Getenv("DB_MAX_CONN_IDLE_TIME"), 30),

		MaxSequencePagination: utils.SafeAtoi(os.Getenv("MAX_SEQUENCE_PAGINATION"), 50),

		CacheLifeWindow: utils.SafeAtoi(os.Getenv("CACHE_LIFE_WINDOW"), 30),
		MaxCacheMemory:  utils.SafeAtoi(os.Getenv("MAX_CACHE_MEMORY"), 10),
	}
}
