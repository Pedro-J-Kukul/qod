// Filename: cmd/api/handlers.go
// Description: HTTP request handlers for the API

package main

import (
	"net/http"
)

// give data, handler function to automate converting to json

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	// create map to hold the json
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     app.config.version,
		},
	}

	// call helper function to write to json
	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The Server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

}
