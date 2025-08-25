package main

import (
	"encoding/json"
	"net/http"
)

// Helper function to write json, has the following parameters:
// response writer, status code for the server, data of type any to encode, and the headers to specify.
// returns an error
func (a *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {

	// encodes data to json
	jsResponse, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// append json and add a new line after each appendage
	jsResponse = append(jsResponse, '\n')

	// set key value pairs
	for key, value := range headers {
		w.Header()[key] = value

		// w.Header().Set*key, value[0]
	}

	// set content type to header

	w.Header().Set("Content-Type", "application/json")

	// write to header
	w.WriteHeader(status)

	// write the json to the body
	w.Write(jsResponse)

	// returns no error/empty
	return nil
}
