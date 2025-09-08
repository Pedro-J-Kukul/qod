package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Function: serve
// Description: Starts the HTTP server
func (app *appDependencies) serve() error {

	// create a new http server with some sensible timeout settings
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	// log the server start
	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)

	// start the server
	return srv.ListenAndServe()
}
