package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthcheck", app.requireAuthentication(app.healthcheck))
	mux.HandleFunc("POST /v1/user/register", app.registerHadler)
	mux.HandleFunc("POST /v1/user/login", app.loginHandler)

	return app.authenticate(mux)
}
