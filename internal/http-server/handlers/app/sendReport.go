package app

import (
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type ReportSender interface {
	SaveReport(data Report) (int64, error)
}
type Report struct {
	Photo        []byte
	TypeOfReport string
	Location     string
	Comment      string
	Time         time.Time
	Phone        string
	Email        string
}

func SendReport(log *slog.Logger, sender ReportSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "app.handlers.sendReport"
		log = log.With(
			slog.String("op", op),
		)
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Error("failed parsing multipart form", slog.String("error", err.Error()))
			render.Status(r, http.StatusNoContent)
			render.JSON(w, r, render.M{})
			return
		}
		data, err := ReadFormData(r)
		if err != nil {
			log.Error("failed reading form data", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, render.M{})
		}
		ID, err := sender.SaveReport(data)
		if err != nil {
			log.Error("failed saving report", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, render.M{})
		}
		render.JSON(w, r, render.M{"ID": ID})
	}
}

func ReadFormData(r *http.Request) (Report, error) {
	const op = "app.readFormData"
	var report Report
	var (
		Photo        []byte
		TypeOfReport string
		Location     string
		Comment      string
		Time         time.Time
		Phone        string
		Email        string
	)
	file, _, err := r.FormFile("photo")
	if err != nil {
		return report, errors.Errorf("%s: failed reading form file", op)
	}
	defer file.Close()
	Photo, _ = io.ReadAll(file)
	TypeOfReport = r.FormValue("type")
	Location = r.FormValue("location")
	Comment = r.FormValue("comment")
	i, err := strconv.ParseInt(r.FormValue("time"), 10, 64)
	if err != nil {
		return report, errors.Errorf("%s: failed parsing time: %s", op, err.Error())
	}
	Time = time.Unix(i, 0)
	Phone = r.FormValue("phone")
	Email = r.FormValue("email")
	report.Photo = Photo
	report.TypeOfReport = TypeOfReport
	report.Location = Location
	report.Comment = Comment
	report.Time = Time
	report.Phone = Phone
	report.Email = Email
	return report, nil
}
