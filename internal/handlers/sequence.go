package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/murilo-bracero/sequence-technical-test/internal/dto"
	"github.com/murilo-bracero/sequence-technical-test/internal/server/cache"
	"github.com/murilo-bracero/sequence-technical-test/internal/server/config"
	"github.com/murilo-bracero/sequence-technical-test/internal/services"
	"github.com/murilo-bracero/sequence-technical-test/internal/utils"
)

type SequenceHandler interface {
	GetSequences(w http.ResponseWriter, r *http.Request)
	GetSequence(w http.ResponseWriter, r *http.Request)
	UpdateSequence(w http.ResponseWriter, r *http.Request)
	CreateSequence(w http.ResponseWriter, r *http.Request)
}

type sequenceHandler struct {
	cfg             *config.Config
	cache           cache.Cache
	sequenceService services.SequenceService
}

func NewSequenceHandler(cfg *config.Config, cache cache.Cache, sequenceService services.SequenceService) *sequenceHandler {
	return &sequenceHandler{cfg: cfg, cache: cache, sequenceService: sequenceService}
}

func (h *sequenceHandler) GetSequences(w http.ResponseWriter, r *http.Request) {
	size := utils.SafeAtoi(r.URL.Query().Get("size"), 50)

	size = min(size, h.cfg.MaxSequencePagination)

	page := utils.SafeAtoi(r.URL.Query().Get("page"), 0)

	key := fmt.Sprintf("sequences-%d-%d", size, page)

	if h.cache.Get(key) != nil {
		w.Write(h.cache.Get(key))
		return
	}

	sequences, err := h.sequenceService.GetSequences(r.Context(), size, page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(sequences)

	if raw, err := json.Marshal(sequences); err == nil {
		h.cache.Set(key, raw)
	}
}

func (h *sequenceHandler) GetSequence(w http.ResponseWriter, r *http.Request) {
	if h.cache.Get("sequence-"+r.PathValue("id")) != nil {
		w.Write(h.cache.Get("sequence-" + r.PathValue("id")))
		return
	}

	id := r.PathValue("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sequence, err := h.sequenceService.GetSequence(r.Context(), uid)
	if err != nil {
		if err == services.ErrorSequenceNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(sequence)

	if raw, err := json.Marshal(sequence); err == nil {
		h.cache.Set("sequence-"+id, raw)
	}
}

func (h *sequenceHandler) UpdateSequence(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req dto.UpdateSequenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sequence, err := h.sequenceService.UpdateSequence(r.Context(), uid, req)
	if err != nil {
		if err == services.ErrorSequenceNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sequence)

	h.cache.EvictAll()
}

func (h *sequenceHandler) CreateSequence(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSequenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&dto.HTTPError{Message: err.Error()})
		return
	}

	sequence, err := h.sequenceService.CreateSequence(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sequence)

	h.cache.EvictAll()
}
