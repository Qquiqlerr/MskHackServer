package app

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type ResponseGetAllReports struct {
	Reports []ReportCutted `json:"reports"`
}
type ReportCutted struct {
	StatusID     string    `json:"status_id"`
	StatusName   string    `json:"status_name"`
	TypeOfReport string    `json:"type_of_report"`
	TypeID       int       `json:"type_id"`
	Coords       []float64 `json:"coords"`
}
type ReportsGetter interface {
	GetAllReports() ([]ReportCutted, error)
}

func GetAllReports(log *slog.Logger, getter ReportsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.app.getAllReports"
		log = log.With(slog.String("operation", op))
		reps, err := getter.GetAllReports()
		if err != nil {
			log.Error(err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, render.M{})
		}
		render.JSON(w, r, ResponseGetAllReports{reps})
	}
}
