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

	cont := controllers.NewControllers(DB)

	fileServer := http.FileServer(http.Dir(filepathRoot))
	mux := http.NewServeMux()
	mux.Handle("/app/*", cont.ApiOps.MiddlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /admin/metrics", cont.ApiOps.AdminMetricsHandler)
	mux.HandleFunc("GET /api/healthz", controllers.ReadinessHandler)
	mux.HandleFunc("GET /api/reset", cont.ApiOps.ResetHits)
	mux.HandleFunc("POST /api/chirps", cont.Chirps.CreateChirpHandler)
	mux.HandleFunc("GET /api/chirps", cont.Chirps.GetChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cont.Chirps.GetSingleChirpHandler)
	mux.HandleFunc("POST /api/users", cont.Users.CreateUserHandler)
	mux.HandleFunc("POST /api/login", cont.Users.LoginHandler)

	// Setup and run server
	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Starting sever on port %s...\n", port)
	log.Fatal(srv.ListenAndServe())
}
