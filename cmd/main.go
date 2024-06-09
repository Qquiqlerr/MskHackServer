package main

import (
	"github.com/go-chi/chi/v5"
	"greenkmchSever/internal/config"
	"greenkmchSever/internal/http-server/handlers/app"
	"greenkmchSever/internal/storage/postgres"
	"log/slog"
	"net/http"
	"os"
)

func main() {

	cfg := config.MustLoad()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	storage, err := postgres.New(cfg.StorageURL)
	if err != nil {
		log.Error("Failed to initialize storage", err.Error())
		os.Exit(1)
	}
	router := chi.NewRouter()
	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/static/*", http.StripPrefix("/static", fs))

	router.Route("/api/mapp", func(r chi.Router) {
		r.Get("/routes_all", app.RoutesAll(log, storage, cfg.Addr))
		r.Post("/visit_request", app.SendVisitRequest(log, storage))
		r.Get("/visit_request", app.GetVisitRequestStatus(log, storage))
		r.Post("/send_report", app.SendReport(log, storage))
		r.Get("/get_all_reports", app.GetAllReports(log, storage))
	})

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start server", err.Error())
		os.Exit(1)
	}
	log.Info("Server shutdown")
}
