// Filename: cmd/api/handlers.go
// Description: General handlers for the API

package main

import (
	"net/http"
)

// When sending a response, we send the header first and then the body
// give data, handler function to automate converting to json
func (app *appDependencies) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	// create map to hold the json
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     AppVersion,
		},
	}

	// call helper function to write to json
	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
