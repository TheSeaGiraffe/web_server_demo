package controllers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/TheSeaGiraffe/web_server_demo/internal/models"
)

func (app *Application) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Config.ServerHitsIncrement()
		next.ServeHTTP(w, r)
	})
}

func (app *Application) isTokenJWT(token string) (bool, error) {
	// Attempt to split the token by '.'
	jwtHeader := strings.Split(token, ".")[0]
	decodedJWTHeader, err := base64.StdEncoding.DecodeString(jwtHeader)
	if err != nil {
		return false, err
	}

	// Attempt to unmarshal the header
	var headerJSON struct {
		Algorithm string `json:"alg"`
		Type      string `json:"typ"`
	}
	err = json.Unmarshal([]byte(decodedJWTHeader), &headerJSON)
	if err != nil {
		return false, err
	}

	// Check that the header has "typ" and "alg" fields
	// Need to add a proper check for the "alg" field later
	switch {
	case (headerJSON.Algorithm == "") || (headerJSON.Type == ""):
		return false, nil
	case headerJSON.Type != "JWT":
		return false, nil
	default:
		return true, nil
	}
}

func (app *Application) MiddlewareAuthenticateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidCredentialsResponse(w, r)
			return
		}

		token := headerParts[1]
		ok, err := app.isTokenJWT(token)
		if err != nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Should be a refresh token so we don't authenticate the user
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		userID, err := app.getIDFromJWT(token)
		if err != nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.DB.GetUserByID(userID)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrUserNotExist):
				app.invalidAuthenticationTokenResponse(w, r)
				return
			default:
				app.serverErrorResponse(w, r)
				return
			}
		}

		r = app.contextSetUser(r, &user)

		next.ServeHTTP(w, r)
	})
}

func (app *Application) MiddlewareAuthenticateRefresh(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if user == nil {
			app.authenticationRequiredResponse(w, r)
			return
		}
		next(w, r)
	}
}

func (app *Application) MiddlewareRequireUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if user == nil {
			app.authenticationRequiredResponse(w, r)
			return
		}
		next(w, r)
	}
}
