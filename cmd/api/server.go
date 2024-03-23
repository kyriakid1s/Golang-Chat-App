package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (app *application) serve() error {

	server := http.Server{
		Addr:     fmt.Sprintf(":%d", app.config.port),
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	app.logger.Info("Starting Server...", "addr", server.Addr, "env", app.config.env)
	return server.ListenAndServe()
}
