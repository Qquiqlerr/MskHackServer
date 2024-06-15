package portal

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type OOPTStress struct {
	Stress float64 `json:"stress"`
	ZoneID int     `json:"zone_id"`
}

type OOPTStressSender interface {
	SendOOPTStress(stress float64, ID int) error
}

func SendOOPTStress(log *slog.Logger, sender OOPTStressSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.portal.sendOOPTStress"
		log = log.With(slog.String("operation", op))
		var stress OOPTStress
		err := render.DecodeJSON(r.Body, &stress)
		if err != nil {
			log.Error("failed to parse stress", err)
			render.Status(r, http.StatusBadRequest)
			return
		}
		err = sender.SendOOPTStress(stress.Stress, stress.ZoneID)
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			return
		}
		render.Status(r, http.StatusOK)
	}
}
