package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"

	"github.com/TheSeaGiraffe/web_server_demo/internal/models"
)

type UsersController struct {
	DB *models.DB
}

func (u *UsersController) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON from the response body
	var input struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// Check that the email address is valid
	_, err = mail.ParseAddress(input.Email)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Not a valid email")
		return

	}

	// Create user
	user, err := u.DB.CreateUser(input.Email)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	// Response is valid
	err = writeJSON(w, http.StatusCreated, user, nil)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
	}
}
