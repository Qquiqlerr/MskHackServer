package app

import (
	"fmt"
	"github.com/go-chi/render"
	"greenkmchSever/internal/lib/files"
	"log/slog"
	"net/http"
)

type Response struct {
	Routes []Routes `json:"routes"`
}
type Routes struct {
	ID       int      `json:"id" db:"id"`
	Zone_id  int      `json:"zone_id" db:"zone_id"`
	Name     string   `json:"name" db:"name"`
	Desc     string   `json:"desc" db:"description"`
	Duration int      `json:"duration" db:"lenght_hours"`
	Img_urls []string `json:"img_urls"`
}
type RoutesGetter interface {
	GetAllRoutes() ([]Routes, error)
}

func RoutesAll(log *slog.Logger, getter RoutesGetter, address string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "app.handlers.routesAll"
		log = log.With(
			slog.String("operation", op),
		)
		var resp Response
		routes, err := getter.GetAllRoutes()
		if err != nil {
			log.Warn("failed to retrieve routes", err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp)
			return
		}
		resp.Routes = routes
		for i, route := range routes {
			filePath := fmt.Sprintf("./static/appdata/%d/%d/imgs", route.Zone_id, route.ID)
			filesList, err := files.GetFilesList(filePath, address)
			if err != nil {
				log.Warn("failed to retrieve files list from file path", filePath)
			}
			routes[i].Img_urls = filesList
		}
		render.JSON(w, r, resp)
	}
}
