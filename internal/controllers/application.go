package controllers

import (
	"fmt"
	"net/http"

	"github.com/TheSeaGiraffe/web_server_demo/internal/models"
)

type Application struct {
	Config ApiConfig
	DB     *models.DB
}

func (app *Application) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (app *Application) ResetHitsHandler(w http.ResponseWriter, r *http.Request) {
	app.Config.ServerHitsReset()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File server hits count has been reset to 0"))
}

func (app *Application) AdminMetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
<html>
    <body>
        <h1>Welcome, Chirpy Admin</h1>
        <p>Chirpy has been visited %d times!</p>
    </body>
</html>`, app.Config.ServerHitsGet())))
}
