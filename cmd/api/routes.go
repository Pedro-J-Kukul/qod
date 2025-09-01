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

	// Handle 404 errors
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// handling 405 errors
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// handler for the healthcheck api
	router.HandlerFunc(http.MethodGet, app.healthCheckName(), app.healthcheckHandler)

	// return router to call appropriate handlers
	// include panic middleware
	return app.recoverPanic(router)
}
