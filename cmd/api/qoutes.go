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

// Handler for retrieving a qoute
func (a *appDependencies) displayQouteHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id from the URL
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Call the Get method on the model to retrieve the data
	quote, err := a.model.Get(id)
	if err != nil {
		switch {
		case err == data.ErrRecordNotFound:
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// display the quote
	data := envelope{"qoute": quote}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

// Update Qoute Handler
func (a *appDependencies) updateQouteHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id from the URL
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Call the Get method on the model to retrieve the data
	quote, err := a.model.Get(id)
	if err != nil {
		switch {
		case err == data.ErrRecordNotFound:
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// create a struct to hold the incoming qoute data
	var incomingData struct {
		Type   *string `json:"type"`
		Quote  *string `json:"quote"`
		Author *string `json:"author"`
	}

	// Perform the decoding
	err = a.readJson(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		// return to prevent further processing
		return
	}

	// check which fields were provided and update the record accordingly
	if incomingData.Type != nil {
		quote.Type = *incomingData.Type
	}
	if incomingData.Quote != nil {
		quote.Qoute = *incomingData.Quote
	}
	if incomingData.Author != nil {
		quote.Author = *incomingData.Author
	}

	// validator instance
	v := validator.NewValidator()

	// validate the incoming data
	data.ValidateQoute(v, quote)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.model.Update(quote)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// send a JSON response with 200 status code
	data := envelope{"qoute": quote}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

// Delete Qoute Handler
func (a *appDependencies) deleteQouteHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id from the URL
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.model.Delete(id)
	if err != nil {
		switch {
		case err == data.ErrRecordNotFound:
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{"message": "qoute successfully deleted"}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *appDependencies) listQoutesHandler(w http.ResponseWriter, r *http.Request) {

	// create a struct to hold the query string parameters
	// type, author, page, page_size
	var queryParamterData struct {
		Type   string
		Quote  string
		Author string
	}

	queryParamters := r.URL.Query()

	// read the values from the query string into the struct
	queryParamterData.Type = a.getSingleQueryParam(queryParamters, "type", "")
	queryParamterData.Quote = a.getSingleQueryParam(queryParamters, "quote", "")
	queryParamterData.Author = a.getSingleQueryParam(queryParamters, "author", "")

	qoutes, err := a.model.GetAll(queryParamterData.Type, queryParamterData.Quote, queryParamterData.Author)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// send a JSON response with 200 status code
	data := envelope{"qoutes": qoutes}
	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}
