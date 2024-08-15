package main

import (
	"learn_web_servers/web_server_demo/internal/handlers"
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	// Setup routing
	apiCfg := handlers.NewApiConfig()
	fileServer := http.FileServer(http.Dir(filepathRoot))
	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /admin/metrics", apiCfg.AdminMetricsHandler)
	mux.HandleFunc("GET /api/healthz", handlers.ReadinessHandler)
	mux.HandleFunc("GET /api/reset", apiCfg.ResetHits)
	mux.HandleFunc("POST /api/validate_chirp", handlers.ValidateChirpHandler)

	// Setup and run server
	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Starting sever on port %s...\n", port)
	log.Fatal(srv.ListenAndServe())
}
