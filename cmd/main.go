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
