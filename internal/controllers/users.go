package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"

	"github.com/TheSeaGiraffe/web_server_demo/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Implement a better validation system later. For now, just make sure that everything works.

func (app *Application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON from the response body
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// Check that the email address is valid
	_, err = mail.ParseAddress(input.Email)
	if err != nil {
		app.errorResponse(w, http.StatusBadRequest, "Not a valid email")
		return

	}

	// Check that email address is not being used
	// Figure out a better way of doing this later
	exists, err := app.DB.EmailExists(input.Email)
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, "Problems accessing database")
		return
	}

	if exists {
		app.errorResponse(w, http.StatusBadRequest, "Account with that email address already exists")
		return
	}

	// Check that the password is valid
	// Will take out the more stringent password requirements for now
	if input.Password == "" {
		app.errorResponse(w, http.StatusBadRequest, "Password must be provided")
		return
	}
	// if len(input.Password) < 8 {
	// 	errorResponse(w, http.StatusBadRequest, "Password must be at least 8 characters long.")
	// 	return
	// }
	// if len(input.Password) > 72 {
	// 	errorResponse(w, http.StatusBadRequest, "Password must not be more than 72 characters long.")
	// 	return
	// }

	// Create user
	user, err := app.DB.CreateUser(input.Email, input.Password)
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	// Response is valid
	output := struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}{
		user.ID,
		user.Email,
	}
	err = app.writeJSON(w, http.StatusCreated, output, nil)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
	}
}

func (app *Application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON from the response body
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// Check if user exists in database and retrieve their info if they do
	user, err := app.DB.GetUserByEmail(input.Email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotExist) {
			app.errorResponse(w, http.StatusBadRequest, "No user with this email exists")
			return
		}
		app.errorResponse(w, http.StatusInternalServerError, "Problem checking database")
		return
	}

	// Compare the password in the request to the existing user's password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		app.errorResponse(w, http.StatusUnauthorized, "Password is incorrect")
		return
	}

	// Return user info sans password on successful login
	output := struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}{
		user.ID,
		user.Email,
	}
	err = app.writeJSON(w, http.StatusOK, output, nil)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
	}
}
