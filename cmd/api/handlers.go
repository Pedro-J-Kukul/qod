// Filename: cmd/api/handlers.go
// Description: HTTP request handlers for the API

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// give data, handler function to automate converting to json

func (app *application) v2healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	// create map to hold the json
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     app.config.version,
	}

	// using json.Marshal, encode the data to json
	jsResponse, err := json.Marshal(data)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request.", http.StatusInternalServerError)

		return
	}

	jsResponse = append(jsResponse, '\n')

	// set header
	w.Header().Set("Content-Type", "application/json")

	// write the json to the body
	w.Write(jsResponse)

}

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	js := `{"status": "available", 
	"environment": %q, 
	"version": %q
	}`
	js = fmt.Sprintf(js, app.config.env, app.config.version)

	// Content-Type is text/plain by default

	// Always set the content type as a json so that the handler sends back json
	w.Header().Set("Content-Type", "application/json")
	// Write the JSON as the HTTP response body.
	w.Write([]byte(js))

}
