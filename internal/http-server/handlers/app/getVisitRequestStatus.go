package app

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type StatusGetter interface {
	GetStatus(fn, sn, mn string) (string, error)
}
type RequestStatus struct {
	Status string `json:"status"`
}

func GetVisitRequestStatus(log *slog.Logger, getter StatusGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.app.GetVisitRequestStatus"
		log = log.With("operation", op)
		var fn, sn, mn string
		var status RequestStatus
		params := r.URL.Query()
		fn = params.Get("fn")
		sn = params.Get("ln")
		mn = params.Get("mn")
		st, err := getter.GetStatus(fn, sn, mn)
		if err != nil {
			log.Warn("Failed to get status", "fn", fn, "sn", sn, "mn", mn, "err", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, status)
			return
		}
		status.Status = st
		render.JSON(w, r, status)
	}
}
