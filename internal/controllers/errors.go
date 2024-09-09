package controllers

import "net/http"

func (app *Application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.errorResponse(w, http.StatusUnauthorized, message)
}

func (app *Application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	app.errorResponse(w, http.StatusUnauthorized, message)
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request) {
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, http.StatusInternalServerError, message)
}

func (app *Application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.errorResponse(w, http.StatusUnauthorized, message)
}
