package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/murilo-bracero/sequence-technical-test/internal/utils"
)

type Config struct {
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	MaxDbConnections int
	MinDbConnections int
	MaxConnIdleTime  int

	MaxSequencePagination int
}

func New() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Warn("failed to load .env file", err.Error(), err)
	}

	return &Config{
		PostgresHost:     os.Getenv("DB_HOST"),
		PostgresPort:     utils.SafeAtoi(os.Getenv("DB_PORT"), 5432),
		PostgresUser:     os.Getenv("DB_USER"),
		PostgresPassword: os.Getenv("DB_PASSWORD"),
		PostgresDatabase: os.Getenv("DB_NAME"),
		MaxDbConnections: utils.SafeAtoi(os.Getenv("DB_MAX_CONNECTIONS"), 10),
		MinDbConnections: utils.SafeAtoi(os.Getenv("DB_MIN_CONNECTIONS"), 1),
		MaxConnIdleTime:  utils.SafeAtoi(os.Getenv("DB_MAX_CONN_IDLE_TIME"), 30),

		MaxSequencePagination: utils.SafeAtoi(os.Getenv("MAX_SEQUENCE_PAGINATION"), 50),
	}
}
