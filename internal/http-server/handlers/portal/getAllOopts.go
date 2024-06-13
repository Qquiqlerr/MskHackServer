package portal

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type OoptsGetter interface {
	GetAllOopts() ([]Oopt, error)
}

type Response struct {
	Oopts []Oopt `json:"oopts"`
}
type Oopt struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

func GetAllOopts(log *slog.Logger, getter OoptsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "portal.GetAllOopts"
		log = log.With(
			slog.String("op", op),
		)
		oopts, err := getter.GetAllOopts()
		if err != nil {
			log.Error("failed to get oopts", err)
			render.Status(r, http.StatusInternalServerError)
			return
		}
		render.JSON(w, r, Response{Oopts: oopts})
	}
}
