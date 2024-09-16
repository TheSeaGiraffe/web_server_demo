package controllers

import (
	"log"
	"net/http"
)

func (app *Application) RefreshAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Get user from request context
	user := app.contextGetUser(r)

	// Create new JWT for the current user
	token, err := app.generateJWT(user.ID, nil)
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, "Could not create JWT")
		return
	}

	// Return newly created JWT in response
	output := struct {
		Token string `json:"token"`
	}{
		token,
	}
	err = app.writeJSON(w, http.StatusOK, output, nil)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
	}
}
