package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/murilo-bracero/sequence-technical-test/internal/dto"
	"github.com/murilo-bracero/sequence-technical-test/internal/server/cache"
	"github.com/murilo-bracero/sequence-technical-test/internal/services"
)

type StepHandler interface {
	CreateStep(w http.ResponseWriter, r *http.Request)
	UpdateStep(w http.ResponseWriter, r *http.Request)
	DeleteStep(w http.ResponseWriter, r *http.Request)
}

type stepHandler struct {
	cache       cache.Cache
	stepService services.StepService
}

var _ StepHandler = (*stepHandler)(nil)

func NewStepHandler(cache cache.Cache, stepService services.StepService) *stepHandler {
	return &stepHandler{stepService: stepService, cache: cache}
}

func (h *stepHandler) CreateStep(w http.ResponseWriter, r *http.Request) {
	sequenceID, err := uuid.Parse(r.PathValue("sequence_id"))
	if err != nil {
		slog.Warn("failed to parse sequence id", err.Error(), err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req dto.CreateStepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	step, err := h.stepService.CreateStep(r.Context(), sequenceID, req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(step)

	h.cache.Evict("sequence-" + r.PathValue("sequence_id"))
}

func (h *stepHandler) UpdateStep(w http.ResponseWriter, r *http.Request) {
	sequenceId := r.PathValue("sequence_id")
	stepId := r.PathValue("step_id")

	seqid, err := uuid.Parse(sequenceId)
	if err != nil {
		slog.Warn("failed to parse sequence id", err.Error(), err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	stid, err := uuid.Parse(stepId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req dto.UpdateStepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	step, err := h.stepService.UpdateStep(r.Context(), seqid, stid, req)
	if err != nil {
		if err == services.ErrorStepNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(step)

	h.cache.Evict("sequence-" + sequenceId)
}

func (h *stepHandler) DeleteStep(w http.ResponseWriter, r *http.Request) {
	sequenceId := r.PathValue("sequence_id")
	stepId := r.PathValue("step_id")

	_, err := uuid.Parse(sequenceId)
	if err != nil {
		slog.Warn("failed to parse sequence id", err.Error(), err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	stid, err := uuid.Parse(stepId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.stepService.DeleteStep(context.Background(), stid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	h.cache.Evict("sequence-" + sequenceId)
}
