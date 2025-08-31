package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/murilo-bracero/sequence-technical-test/internal/db"
	"github.com/murilo-bracero/sequence-technical-test/internal/handlers"
	"github.com/murilo-bracero/sequence-technical-test/internal/repository"
	"github.com/murilo-bracero/sequence-technical-test/internal/server"
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

	sequenceHandler := handlers.NewSequenceHandler(cfg, sequenceService)

	stepRepository := repository.NewStepRepository(db)

	stepService := services.NewStepService(sequenceRepository, stepRepository)

	stepHandler := handlers.NewStepHandler(stepService)

	if err := server.Start(db, sequenceHandler, stepHandler); err != nil {
		os.Exit(1)
	}
}
