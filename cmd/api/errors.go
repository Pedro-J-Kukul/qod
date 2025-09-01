// Filename: cmd/api/errors.go
// Description: Handling errors in proper json format

package main

import (
	"fmt"
	"net/http"
)

// error helper function to correctly log an error
func (app *application) logError(r *http.Request, err error) {
	method := r.Method
	uri := r.URL.RequestURI()
	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

// Sends an error response in JSON format
func (app *application) errorResponseJSON(w http.ResponseWriter, r *http.Request, status int, message any) {
	// create an envelope of error data
	errorData := envelope{"error": message}
	err := app.writeJSON(w, status, errorData, nil)
	if err != nil {
		// log the error
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// sends an error in case our server is kabloey
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// log the error
	app.logError(r, err)
	// prepare a message to send to the clietn
	message := "the server encountered a problem and could not process your request"
	app.errorResponseJSON(w, r, http.StatusInternalServerError, message)
}

// send an error response if our client messes up with a 404
func (app *application) notFoundResponse(w http.ResponseWriter,
	r *http.Request) {

	// we only log server errors, not client errors
	// prepare a response to send to the client
	message := "the requested resource could not be found"
	app.errorResponseJSON(w, r, http.StatusNotFound, message)
}

// send an error response if our client messes up with a 405
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {

	// we only log server errors, not client errors
	// prepare a formatted response to send to the client
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)

	app.errorResponseJSON(w, r, http.StatusMethodNotAllowed, message)
}

// send an error response if our client messes up with a 400 (bad request)
func (a *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {

	a.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
}
