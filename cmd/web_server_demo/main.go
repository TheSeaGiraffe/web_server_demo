package main

import (
	"log"
	"net/http"
	"os"

	"github.com/TheSeaGiraffe/web_server_demo/internal/controllers"
	"github.com/TheSeaGiraffe/web_server_demo/internal/models"
	"github.com/joho/godotenv"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	// Init DB connection
	DB, err := models.NewDB(models.DBFilePath)
	if err != nil {
		log.Fatalf("Could not connect to DB: %s", err)
	}

	// Init config
	// Maybe think about putting this in a separate function or even package
	// I think it's okay to keep it like this for now
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Could not load environment variables: %s", err)
	}
	jwtSecret := os.Getenv("JWT_SECRET")

	cfg := controllers.NewApiConfig(jwtSecret)

	// Setup the routes
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
	mux.HandleFunc("PUT /api/users", application.MiddlewareRequireUser(application.UpdateUserHandler))
	mux.HandleFunc("POST /api/login", application.LoginHandler)
	mux.HandleFunc("POST /api/refresh",
		application.MiddlewareAuthenticateRefresh(application.MiddlewareRequireUser(application.RefreshAccessTokenHandler)))
	mux.HandleFunc("POST /api/revoke",
		application.MiddlewareAuthenticateRefresh(application.MiddlewareRequireUser(application.RevokeRefreshTokenHandler)))

	// Setup and run server
	srv := http.Server{
		Addr:    ":" + port,
		Handler: application.MiddlewareAuthenticateJWT(mux), // Find a better way of doing this
	}
	log.Printf("Starting sever on port %s...\n", port)
	log.Fatal(srv.ListenAndServe())
}
