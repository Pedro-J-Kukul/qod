// Filename: internal/data/qoutes.go
// Description Handler for the POST method of sending a qoute

package main

import (
	"fmt"
	"net/http"

	"github.com/Pedro-J-Kukul/qod/internal/data"
	"github.com/Pedro-J-Kukul/qod/internal/validator"
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
	// create a qoute struct to hold the data for insertion
	quote := &data.Qoute{
		Type:   incomingData.Type,
		Qoute:  incomingData.Quote,
		Author: incomingData.Author,
	}

	// validator instance
	v := validator.NewValidator()

	// validate the incoming data
	data.ValidateQoute(v, quote)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.model.Insert(quote)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/qoutes/%d", quote.ID))

	// send a JSON response with 201 status code
	data := envelope{"qoute": quote}
	err = a.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}
