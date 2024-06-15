package portal

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Data struct {
	ID     int `json:"id"`
	Status int `json:"status"`
}
type RequestStatusUpdater interface {
	UpdateRequestStatus(id int, status int) error
}

func ChangeRequestStatus(log *slog.Logger, updater RequestStatusUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "portal.handlers.ChangeRequestStatus"
		log = log.With(
			slog.String("operation", op),
		)
		var data Data
		err := render.DecodeJSON(r.Body, &data)
		if err != nil {
			log.Error("failed to bind data", err.Error())
			render.Status(r, http.StatusBadRequest)
			return
		}
		err = updater.UpdateRequestStatus(data.ID, data.Status)
		if err != nil {
			log.Error("failed to update request status", err.Error())
			render.Status(r, http.StatusInternalServerError)
			return
		}
	}
}
