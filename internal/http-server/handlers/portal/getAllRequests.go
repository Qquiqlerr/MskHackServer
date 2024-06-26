package portal

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type ResponseRequest struct {
	Requests []ZonesForDate `json:"requests"`
}

type ZonesForDate struct {
	ID             int64   `json:"id"`
	Route          string  `json:"route"`
	Stress         float64 `json:"stress"`
	StressIfSubmit float64 `json:"stress_if_submit"`
}

type ZonesGetter interface {
	GetAllZonesForDate() ([]ZonesForDate, error)
}

func GetAllRequests(log *slog.Logger, zonesGetter ZonesGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "portal.GetAllRequests"
		log := log.With(slog.String("op", op))

		requests, err := zonesGetter.GetAllZonesForDate()
		log.Info("get zones", "err", err, "requests", requests)
		if err != nil {
			log.Error("failed to get zones", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, ResponseRequest{})
			return
		}

		render.JSON(w, r, ResponseRequest{
			Requests: requests,
		})
	}
}
