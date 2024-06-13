package portal

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	ID        int64  `json:"id"`
	NewStatus string `json:"new_status"`
}

type ProblemUpdater interface {
	UpdateProblem(id int64, newStatus string) error
}

// PutProblem creates an HTTP handler function for updating a problem.
//
// Takes a logger and a ProblemUpdater interface as parameters.
// Returns an http.HandlerFunc.
func PutProblem(log *slog.Logger, updater ProblemUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "portal.PutProblem"
		log = log.With(
			slog.String("op", op),
		)
		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", err)
			render.Status(r, http.StatusBadRequest)
			return
		}
		if err := updater.UpdateProblem(req.ID, req.NewStatus); err != nil {
			log.Error("failed to update problem", err)
			render.Status(r, http.StatusInternalServerError)
			return
		}
		render.Status(r, http.StatusOK)
	}
}
