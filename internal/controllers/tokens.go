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

func (app *Application) RevokeRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Get user from request context
	user := app.contextGetUser(r)

	// Get token by user id
	token, err := app.DB.GetTokenByUserID(user.ID)
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, "The user does not have a refresh token")
		return
	}

	// Delete the token using the token's id
	err = app.DB.DeleteRefreshToken(token.ID)
	if err != nil {
		app.serverErrorResponse(w, r)
		return
	}

	// Only send error code 204
	w.WriteHeader(http.StatusNoContent)
}
