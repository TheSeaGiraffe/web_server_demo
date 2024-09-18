package controllers

import (
	"cmp"
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/TheSeaGiraffe/web_server_demo/internal/models"
)

const replacementString = "****"

var badWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

func replaceBadWords(chirp string) string {
	chirpSplit := strings.Fields(chirp)
	for i, word := range chirpSplit {
		_, ok := badWords[strings.ToLower(word)]
		if ok {
			chirpSplit[i] = replacementString
		}
	}
	return strings.Join(chirpSplit, " ")
}

func (app *Application) CreateChirpHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON from the response body
	var input struct {
		Body string `json:"body"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// Check that response is <= 140 chars
	const maxChirpLength = 140
	if len(input.Body) > maxChirpLength {
		app.errorResponse(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// Remove profanity and get user ID
	userID := (app.contextGetUser(r)).ID
	cleanedBody := replaceBadWords(input.Body)
	chirp, err := app.DB.CreateChirp(cleanedBody, userID)
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	// Response is valid
	// output := map[string]string{"cleaned_body": cleanedBody}
	// err = writeJSON(w, http.StatusOK, output, nil)
	err = app.writeJSON(w, http.StatusCreated, chirp, nil)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
	}
}

func (app *Application) GetChirpsHandler(w http.ResponseWriter, r *http.Request) {
	// Get chirps from database
	chirps, err := app.DB.GetChirps()
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, "Couldn't load chirps from database")
		return
	}

	// Sort the chirps and then filter by "author_id" if it is provided
	var chirpsFiltered []models.Chirp
	if len(chirps) > 0 {
		slices.SortFunc(chirps, func(a, b models.Chirp) int {
			return cmp.Compare(a.ID, b.ID)
		})

		var authorID int
		authorIDString := r.URL.Query().Get("author_id")
		if authorIDString != "" {
			authorID, err = strconv.Atoi(authorIDString)
			if err != nil {
				app.serverErrorResponse(w, r)
				return
			}
			for _, chirp := range chirps {
				if chirp.AuthorID == authorID {
					chirpsFiltered = append(chirpsFiltered, chirp)
				}
			}
		} else {
			chirpsFiltered = chirps
		}

	}

	// Return the chirps in a json response
	err = app.writeJSON(w, http.StatusOK, chirpsFiltered, nil)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
	}
}

func (app *Application) GetSingleChirpHandler(w http.ResponseWriter, r *http.Request) {
	// Get chirp ID from URL path
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDStr)
	if err != nil {
		app.errorResponse(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	// Get the chirp with the specified ID
	chirp, err := app.DB.GetChirpByID(chirpID)
	if err != nil {
		app.errorResponse(w, http.StatusNotFound, "Chirp with that ID doesn't exist")
		return
	}

	err = app.writeJSON(w, http.StatusOK, chirp, nil)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
	}
}

func (app *Application) DeleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	// Get chirp ID from URL path
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDStr)
	if err != nil {
		app.errorResponse(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	// Get the chirp with the specified ID
	chirp, err := app.DB.GetChirpByID(chirpID)
	if err != nil {
		app.errorResponse(w, http.StatusNotFound, "Chirp with that ID doesn't exist")
		return
	}

	userID := (app.contextGetUser(r)).ID
	if chirp.AuthorID != userID {
		app.errorResponse(w, http.StatusForbidden, "User is not allowed to access this resource")
		return
	}

	// Delete the chirp with the specified ID
	err = app.DB.DeleteChirpByID(chirpID)
	if err != nil {
		app.errorResponse(w, http.StatusNotFound, "Chirp with that ID doesn't exist")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
