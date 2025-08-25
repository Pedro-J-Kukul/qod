// Filename: cmd/api/handlers.go
// Description: HTTP request handlers for the API

package main

import (
	"fmt"
	"net/http"
	"os"
)

func (app *application) quoteHandler(w http.ResponseWriter, r *http.Request) {
	quotesData, err := os.ReadFile("internal/quotes.json")
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(quotesData)
}

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	js := `{"status": "available", 
	"environment": %q, 
	"version": %q}`
	js = fmt.Sprintf(js, app.config.env, app.config.version)

	// Content-Type is text/plain by default

	w.Header().Set("Content-Type", "application/json")
	// Write the JSON as the HTTP response body.
	w.Write([]byte(js))

}
