package main

import (
	"learn_web_servers/web_server_demo/internal/database"
	"learn_web_servers/web_server_demo/internal/handlers"
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	// Init DB connection
	DB, err := database.NewDB(database.DBFilePath)
	if err != nil {
		log.Fatalf("Could not connect to DB: %s", err)
	}

	// Setup routing
	// Think about creating a unified handler struct
	apiOps := handlers.NewApiOps()
	chirpC := handlers.ChirpController{
		DB: DB,
	}

	fileServer := http.FileServer(http.Dir(filepathRoot))
	mux := http.NewServeMux()
	mux.Handle("/app/*", apiOps.MiddlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /admin/metrics", apiOps.AdminMetricsHandler)
	mux.HandleFunc("GET /api/healthz", handlers.ReadinessHandler)
	mux.HandleFunc("GET /api/reset", apiOps.ResetHits)
	mux.HandleFunc("POST /api/chirps", chirpC.CreateChirpHandler)
	mux.HandleFunc("GET /api/chirps", chirpC.GetChirpsHandler)

	// Setup and run server
	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Starting sever on port %s...\n", port)
	log.Fatal(srv.ListenAndServe())
}
