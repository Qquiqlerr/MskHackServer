package static

import (
	"github.com/go-chi/render"
	"html/template"
	"log/slog"
	"net/http"
)

func RequestInfo(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.BasicAuth()
		ts, err := template.ParseFiles("./static/portal/html/request_info.html")
		if err != nil {
			log.Error("failed to parse template", err.Error())
			render.Status(r, http.StatusInternalServerError)
			return
		}
		err = ts.Execute(w, nil)
		if err != nil {
			log.Error("failed to execute template", err.Error())
			render.Status(r, http.StatusInternalServerError)
			return
		}
	}
}
