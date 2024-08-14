package main

import (
	"learn_web_servers/web_server_demo/internal/apiconfig"
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	// Setup routing
	apiCfg := apiconfig.NewApiConfig()
	fileServer := http.FileServer(http.Dir(filepathRoot))
	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /healthz", readinessHandler)
	mux.HandleFunc("GET /metrics", apiCfg.GetHits)
	mux.HandleFunc("GET /reset", apiCfg.ResetHits)

	// Setup and run server
	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Starting sever on port %s...\n", port)
	log.Fatal(srv.ListenAndServe())
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
