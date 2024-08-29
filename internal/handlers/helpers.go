package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type envelope map[string]any

func writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func errorResponse(w http.ResponseWriter, status int, message any) {
	env := envelope{"error": message}
	err := writeJSON(w, status, env, nil)
	if err != nil {
		log.Printf("Encountered error: %s", err)
		w.WriteHeader(status)
	}
}
