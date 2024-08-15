package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func ValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON from the response body
	var input struct {
		Body string `json:"body"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Couldn't decode parameters")
		w.WriteHeader(500)
		return
	}

	// Check that response is <= 140 chars
	const maxChirpLength = 140
	if len(input.Body) > maxChirpLength {
		errorResponse(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// Response is valid
	output := map[string]bool{"valid": true}
	err = writeJSON(w, http.StatusOK, output, nil)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
	}
}
