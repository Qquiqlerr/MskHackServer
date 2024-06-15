package portal

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type StressSender interface {
	SendRouteStress(int, float64) error
	GetRouteLines(int) ([]float64, error)
}
type Stress struct {
	Route RouteJSON `json:"route"`
}
type RouteJSON struct {
	ID  int     `json:"id"`
	RCC RCCData `json:"RCC"`
}

type RCCData struct {
	PCC PCCData `json:"PCC"`
	MC  float64 `json:"mc"`
}

type PCCData struct {
	BCC BCCData   `json:"BCC"`
	CF  []float64 `json:"cf"`
}

type BCCData struct {
	DG float64 `json:"dg"`
	TS float64 `json:"ts"`
	GS float64 `json:"gs"`
}

func SendRouteStress(log *slog.Logger, storage StressSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "portal.SendRouteStress"

		log = log.With(
			slog.String("op", op),
		)

		var route Stress
		if err := render.DecodeJSON(r.Body, &route); err != nil {
			log.Error("decode json error", err)
			render.Status(r, http.StatusBadRequest)
			return
		}
		stress, err := CalculateStress(route.Route, storage)
		if err != nil {
			log.Error("calculate stress error", err)
			render.Status(r, http.StatusInternalServerError)
			return
		}
		log.Info("send route stress")
		err = storage.SendRouteStress(route.Route.ID, stress)
		if err != nil {
			log.Error("send route stress error", err)
			render.Status(r, http.StatusInternalServerError)
			return
		}
		render.JSON(w, r, "ok")
	}
}
func CalculateStress(route RouteJSON, storage StressSender) (float64, error) {
	var lines []float64
	lines, err := storage.GetRouteLines(route.ID) // get lines from postgres.
	if err != nil {
		return 0, err
	}
	var result float64
	for _, line := range lines {
		result += (line / route.RCC.PCC.BCC.DG) * (route.RCC.PCC.BCC.TS / (line / 3.5))
	}
	result *= route.RCC.PCC.BCC.GS * (1.0 / float64(len(lines)))
	for _, coef := range route.RCC.PCC.CF {
		result *= coef
	}
	result *= route.RCC.MC
	return result, nil
}
