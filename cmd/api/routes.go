// Filename: /cmd/api/routes.go
// Description: connects the routes with an api

package main

import (
	"net/http"

	// Importing Route Package
	"github.com/julienschmidt/httprouter"
)

func (app *appDependencies) routes() http.Handler {

	// create a new router instance
	router := httprouter.New()

	// Handle 404 errors
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// handling 405 errors
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Define the routes
	router.HandlerFunc(http.MethodGet, "/v5/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v8/comments", app.createCommentHandler)
	router.HandlerFunc(http.MethodGet, "/v6/quotes", app.createQouteHandler)

	// include panic middleware
	return app.recoverPanic(router)
}
