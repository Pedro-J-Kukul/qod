// Filename: internal/data/qoutes.go
// Description Handler for the POST method of sending a qoute

package main

import (
	"fmt"
	"net/http"
)

// Handler for creating qoutes
func (a *appDependencies) createQouteHandler(w http.ResponseWriter, r *http.Request) {
	// create a struct to hold the incoming qoute data
	var incomingData struct {
		Type   string `json:"type"`
		Quote  string `json:"quote"`
		Author string `json:"author"`
	}

	// Perform the decoding
	err := a.readJson(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		// return to prevent further processing
		return
	}

	// for now just print the incoming data to the console
	fmt.Fprintf(w, "\n%+v\t\n", incomingData)
}
