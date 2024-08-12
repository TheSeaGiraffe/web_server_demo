package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Setup routing
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.Handle("/assets", http.FileServer(http.Dir("assets")))

	// Setup and run server
	port := "8080"
	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fmt.Printf("Starting sever on port %s...\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
