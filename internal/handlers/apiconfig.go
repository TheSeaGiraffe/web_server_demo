package handlers

import (
	"fmt"
	"net/http"
)

type ApiOps struct {
	fileserverHits int
}

func NewApiOps() *ApiOps {
	return &ApiOps{
		fileserverHits: 0,
	}
}

func (cfg *ApiOps) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiOps) ResetHits(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File server hits count has been reset to 0"))
}

func (cfg *ApiOps) AdminMetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
<html>
    <body>
        <h1>Welcome, Chirpy Admin</h1>
        <p>Chirpy has been visited %d times!</p>
    </body>
</html>`, cfg.fileserverHits)))
}
