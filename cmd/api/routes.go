// Filename: /cmd/api/routes.go
// Description: connects the routes with an api

package main

import (
	"net/http"

	// Importing Route Package
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {

	// create a new router instance
	router := httprouter.New()

	// // handler for the healthcheck api
	// router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// version 2
	router.HandlerFunc(http.MethodGet, "/v4/healthcheck", app.healthcheckHandlerV4)

	// handler for a qoutes api
	// router.HandlerFunc(http.MethodGet, "/v1/quotes", app)

	// return router to call appropriate handlers
	return router
}
