package controllers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/TheSeaGiraffe/web_server_demo/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func (app *Application) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Config.ServerHitsIncrement()
		next.ServeHTTP(w, r)
	})
}

func (app *Application) isTokenJWT(token string) (bool, error) {
	// Attempt to split the token by '.'
	jwtSplit := strings.Split(token, ".")
	if len(jwtSplit) != 3 {
		return false, nil
	}

	jwtHeader := jwtSplit[0]
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

func (app *Application) getIDFromJWT(tokenPlaintext string) (int, error) {
	token, err := jwt.Parse(tokenPlaintext, func(t *jwt.Token) (interface{}, error) {
		return []byte(app.Config.jwtSecret), nil
	})
	// Streamline the error handling logic later
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("Invalid token")
	}

	idStr, err := claims.GetSubject()
	if err != nil {
		return 0, fmt.Errorf("Could not retrieve user ID")
	}

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	return userID, nil
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
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			next(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidCredentialsResponse(w, r)
			return
		}

		// Validate token
		token := headerParts[1]
		tokenLen, err := app.DB.GetRefreshTokenByteLen(token)
		if (err != nil) || (tokenLen != models.RefreshTokenLen) {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Check if token has expired
		isExpired, err := app.DB.RefreshTokenExpired(token)
		if err != nil {
			app.serverErrorResponse(w, r)
			return
		}
		if isExpired {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Get user associated with token
		user, err := app.DB.GetUserByRefreshToken(token)
		if err != nil {
			// In the event of an error, just return the handler
			next(w, r)
			return
		}

		// Add user to context
		r = app.contextSetUser(r, &user)

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
