package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/murilo-bracero/sequence-technical-test/internal/db"
	"github.com/murilo-bracero/sequence-technical-test/internal/handlers"
	"github.com/murilo-bracero/sequence-technical-test/internal/repository"
	"github.com/murilo-bracero/sequence-technical-test/internal/server"
	"github.com/murilo-bracero/sequence-technical-test/internal/server/cache"
	"github.com/murilo-bracero/sequence-technical-test/internal/server/config"
	"github.com/murilo-bracero/sequence-technical-test/internal/services"
)

func main() {
	cfg := config.New()

	db, err := db.New(context.Background(), cfg)
	if err != nil {
		slog.Error("failed to connect to database", err.Error(), err)
		os.Exit(1)
	}

	defer db.Close()

	sequenceRepository := repository.NewSequenceRepository(db)

	sequenceService := services.NewSequenceService(sequenceRepository)

	cache, err := cache.New(context.Background(), cfg)
	if err != nil {
		slog.Error("failed to create cache", err.Error(), err)
		os.Exit(1)
	}

	sequenceHandler := handlers.NewSequenceHandler(cfg, cache, sequenceService)

	stepRepository := repository.NewStepRepository(db)

	stepService := services.NewStepService(sequenceRepository, stepRepository)

	stepHandler := handlers.NewStepHandler(cache, stepService)

	if err := server.Start(cfg, db, sequenceHandler, stepHandler); err != nil {
		os.Exit(1)
	}
}
