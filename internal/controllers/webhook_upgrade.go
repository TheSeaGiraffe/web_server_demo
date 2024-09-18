package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/TheSeaGiraffe/web_server_demo/internal/models"
)

func (app *Application) UpgradeToChirpyRedHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON from the response body
	var input models.EventChirpyRed
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// Check that event is "user.upgraded"
	if input.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Upgrade user
	err = app.DB.UpgradeChirpyRedForUser(input.Data.UserID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotExist) {
			app.errorResponse(w, http.StatusNotFound, "Could not find user")
			return
		}
		app.serverErrorResponse(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
