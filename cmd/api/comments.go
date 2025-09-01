// Filename :cmd/api/comments.go

package main

import (
	"fmt"
	"net/http"
	// internal data
)

func (a *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	// create a struct to thold a comment
	// we use a struct tag [``] to maek the names display in lowercase
	var incomingData struct {
		Content string `json:"content"`
		Author  string `json:"author"`
	}

	// Perform the decoding
	err := a.readJson(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", incomingData)
}
