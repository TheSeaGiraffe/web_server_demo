package main

import (
	"log"
	"net/http"

	"github.com/TheSeaGiraffe/web_server_demo/internal/controllers"
	"github.com/TheSeaGiraffe/web_server_demo/internal/models"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	// Init DB connection
	DB, err := models.NewDB(models.DBFilePath)
	if err != nil {
		log.Fatalf("Could not connect to DB: %s", err)
	}

	cfg := controllers.NewApiConfig()

	application := controllers.Application{
		DB:     DB,
		Config: cfg,
	}

	fileServer := http.FileServer(http.Dir(filepathRoot))
	mux := http.NewServeMux()
	mux.Handle("/app/*", application.MiddlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /admin/metrics", application.AdminMetricsHandler)
	mux.HandleFunc("GET /api/healthz", application.ReadinessHandler)
	mux.HandleFunc("GET /api/reset", application.ResetHitsHandler)
	mux.HandleFunc("POST /api/chirps", application.CreateChirpHandler)
	mux.HandleFunc("GET /api/chirps", application.GetChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", application.GetSingleChirpHandler)
	mux.HandleFunc("POST /api/users", application.CreateUserHandler)
	mux.HandleFunc("POST /api/login", application.LoginHandler)

	// Setup and run server
	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Starting sever on port %s...\n", port)
	log.Fatal(srv.ListenAndServe())
}
