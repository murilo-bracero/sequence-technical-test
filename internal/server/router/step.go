package router

import (
	"net/http"

	"github.com/murilo-bracero/sequence-technical-test/internal/handlers"
)

func StepRouter(stepHandler handlers.StepHandler, r *http.ServeMux) {
	r.HandleFunc("POST /sequences/{sequence_id}/steps", stepHandler.CreateStep)
	r.HandleFunc("PATCH /sequences/{sequence_id}/steps/{step_id}", stepHandler.UpdateStep)
	r.HandleFunc("DELETE /sequences/{sequence_id}/steps/{step_id}", stepHandler.DeleteStep)
}
