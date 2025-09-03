// Filename: internal/data/qoutes.go
// Description Handler for the POST method of sending a qoute

package main

import (
	"fmt"
	"net/http"
)

func (a *application) createQouteHandler(w http.ResponseWriter, r *http.Request) {
	var incomingData struct {
		Type   string `json:"type"`
		Quote  string `json:"quote"`
		Author string `json:"author"`
	}

	err := a.readJson(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Print on a new line the qoute insertedpsql --vers
	fmt.Fprintf(w, "%+v\n", incomingData)
}
