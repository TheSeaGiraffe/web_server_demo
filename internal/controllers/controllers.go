package controllers

import "github.com/TheSeaGiraffe/web_server_demo/internal/models"

type Controllers struct {
	Chirps ChirpController
	Users  UsersController
	ApiOps ApiOps
}

func NewControllers(db *models.DB) Controllers {
	return Controllers{
		ApiOps: ApiOps{fileserverHits: 0},
		Chirps: ChirpController{DB: db},
		Users:  UsersController{DB: db},
	}
}
