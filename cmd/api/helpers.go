package main

import (
	"encoding/json"
	"net/http"
)

// creating an envelope type
type envelope map[string]any

// Helper function to write json, has the following parameters:
// response writer, status code for the server, data of custom type envelope which is a map to encode, and the headers to specify.
// returns an error
func (a *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {

	// encodes data to json
	// use marshall indent to add an indent on each line of json
	jsResponse, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// append json and add a new line after each appendage
	jsResponse = append(jsResponse, '\n')

	// setting addtiional headers
	for key, value := range headers {
		w.Header()[key] = value
	}

	// set content type to header
	w.Header().Set("Content-Type", "application/json")

	// Explicitly setting the response status code
	w.WriteHeader(status)

	// writing the json to the body, but also checking for errors
	_, err = w.Write(jsResponse)
	if err != nil {
		return err
	}

	// returns no error/empty
	return nil
}
