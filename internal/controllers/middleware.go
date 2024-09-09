package controllers

import (
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

func (app *Application) MiddlewareAuthenticate(next http.Handler) http.Handler {
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
