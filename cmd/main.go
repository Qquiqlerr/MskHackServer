package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"greenkmchSever/internal/config"
	"greenkmchSever/internal/http-server/handlers/app"
	"greenkmchSever/internal/http-server/handlers/portal"
	"greenkmchSever/internal/http-server/handlers/portal/static"
	"greenkmchSever/internal/storage/postgres"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// Load the configuration settings
	cfg := config.MustLoad()

	// Initialize logger with debug level
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// Create a new storage instance
	storage, err := postgres.New(cfg.StorageURL)
	if err != nil {
		log.Error("Failed to initialize storage", err.Error())
		os.Exit(1)
	}
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	router := chi.NewRouter()
	router.Use(cors.Handler)
	// Serve static files from the "./static" directory
	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/static/*", http.StripPrefix("/static", fs))

	// Define routes under "/api/mapp"
	router.Route("/api/mapp", func(r chi.Router) {
		// Handle GET request for "/routes_all"
		r.Get("/routes_all", app.RoutesAll(log, storage, cfg.Addr))
		// Handle POST request for "/visit_request"
		r.Post("/visit_request", app.SendVisitRequest(log, storage))
		// Handle GET request for "/visit_request"
		r.Get("/visit_request", app.GetVisitRequestStatus(log, storage))
		// Handle POST request for "/send_report"
		r.Post("/send_report", app.SendReport(log, storage))
		// Handle GET request for "/get_all_reports"
		r.Get("/get_all_reports", app.GetAllReports(log, storage))
	})

	// Define routes under "/api/portal"
	router.Route("/api/portal", func(r chi.Router) {
		// Handle GET request for "/get_all_problems"
		r.Get("/get_all_problems", portal.GetAllProblems(log, storage))
		r.Put("/update_problem", portal.PutProblem(log, storage))
		r.Get("/get_all_oopts", portal.GetAllOopts(log, storage))
		r.Get("/get_all_routes", portal.GetAllRoutesFromZone(log, storage))
		r.Put("/send_route_stress", portal.SendRouteStress(log, storage))
		r.Put("/send_oopt_stress", portal.SendOOPTStress(log, storage))
		r.Get("/get_all_requests", portal.GetAllRequests(log, storage))
	})

	// Define routes under "/portal"
	router.Route("/portal", func(r chi.Router) {
		// Apply basic authentication middleware
		r.Use(middleware.BasicAuth("", map[string]string{
			"admin": "admin",
		}))
		// Handle GET request for "/troubles"
		r.Get("/troubles", static.ListOfTroubles(log))
		// Handle GET request for "/"
		r.Get("/", static.IndexPage(log))
		// Handle GET request for "/list_of_oops"
		r.Get("/list_of_oopts", static.ListOfOopts(log))
		r.Get("/list_of_routes", static.ListOfRoutes(log))
		r.Route("/info", func(r chi.Router) {
			r.Get("/page1", static.Page1(log))
			r.Get("/page2", static.Page2(log))
			r.Get("/page3", static.Page3(log))
		})
		r.Get("/list_of_requests", static.ListOfRequsets(log))
	})

	// Create an HTTP server
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	// Start the server
	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start server", err.Error())
		os.Exit(1)
	}

	// Log server shutdown
	log.Info("Server shutdown")
}
