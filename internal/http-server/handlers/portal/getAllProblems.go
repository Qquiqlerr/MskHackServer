package portal

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Problems struct {
	Problems []Problem `json:"problems"`
}
type Problem struct {
	Id      int64  `json:"Id"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Comment string `json:"comment"`
}

type ProblemsGetter interface {
	GetProblems() (Problems, error)
}

func GetAllProblems(log *slog.Logger, getter ProblemsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "portal.GetAllProblems"
		log = log.With(
			slog.String("op", op),
		)
		problems, err := getter.GetProblems()
		if err != nil {
			log.Error("get all problems error", err)
			render.Status(r, http.StatusInternalServerError)
			return
		}
		render.JSON(w, r, problems)
	}
}
