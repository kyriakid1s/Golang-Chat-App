package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthcheck", app.healthcheck)
	mux.HandleFunc("POST /v1/user/register", app.registerHadler)
	mux.HandleFunc("POST /v1/user/login", app.loginHandler)
	mux.HandleFunc("POST /v1/message/send", app.requireAuthentication(app.sendHandler))
	mux.HandleFunc("GET /v1/message/get/{name}", app.requireAuthentication(app.getHandler))
	mux.HandleFunc("GET /v1/message/getchats", app.requireAuthentication(app.getChatsHandler))
	mux.HandleFunc("GET /v1/ws/getuser", app.requireAuthentication(app.getCurrentUserHandler))
	mux.HandleFunc("/ws", app.handleConnections)

	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(mux))))
}
