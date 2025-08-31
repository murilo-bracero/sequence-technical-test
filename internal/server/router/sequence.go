package router

import (
	"net/http"

	"github.com/murilo-bracero/sequence-technical-test/internal/handlers"
)

func SequenceRouter(sequenceHandler handlers.SequenceHandler, r *http.ServeMux) {
	r.HandleFunc("GET /sequences", sequenceHandler.GetSequences)
	r.HandleFunc("GET /sequences/{id}", sequenceHandler.GetSequence)
	r.HandleFunc("PATCH /sequences/{id}", sequenceHandler.UpdateSequence)
	r.HandleFunc("POST /sequences", sequenceHandler.CreateSequence)
}
