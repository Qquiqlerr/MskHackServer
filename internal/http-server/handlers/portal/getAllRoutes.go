package portal

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type RoutesGetter interface {
	GetAllRoutesFromZone(zoneID int) ([]Route, error)
}
type Route struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Stress float64 `json:"stress"`
}
type Resp struct {
	Routes []Route `json:"routes"`
}

func GetAllRoutesFromZone(log *slog.Logger, getter RoutesGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "portal.GetAllRoutesFromZone"
		log = log.With(
			slog.String("op", op),
		)
		zoneID, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Error("failed to parse zone id", err)
			render.Status(r, http.StatusBadRequest)
			return
		}
		routes, err := getter.GetAllRoutesFromZone(zoneID)
		if err != nil {
			log.Error("failed to get routes", err)
			render.Status(r, http.StatusInternalServerError)
			return
		}
		render.JSON(w, r, Resp{Routes: routes})
	}
}
