package controllers

import "net/http"

func (app *Application) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Config.ServerHitsIncrement()
		next.ServeHTTP(w, r)
	})
}
