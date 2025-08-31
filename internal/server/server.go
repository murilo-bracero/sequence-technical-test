package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/murilo-bracero/sequence-technical-test/internal/db"
	"github.com/murilo-bracero/sequence-technical-test/internal/handlers"
	"github.com/murilo-bracero/sequence-technical-test/internal/server/router"
)

func Start(db db.DB, sequenceHandler handlers.SequenceHandler, stepHandler handlers.StepHandler) error {
	r := http.NewServeMux()

	router.SequenceRouter(sequenceHandler, r)
	router.StepRouter(stepHandler, r)

	r.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		res := make(map[string]string)

		res["app"] = "ok"

		if err := db.Ping(context.Background()); err != nil {
			slog.Error("failed to ping database", err.Error(), err)
			res["database"] = "error"
		} else {
			res["database"] = "ok"
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	slog.Info("Starting server on port 8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
		slog.Error("failed to start server", err.Error(), err)
		return err
	}

	return nil
}
