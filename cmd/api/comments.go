// Filename :cmd/api/comments.go

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	// internal data
)

func (a *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	// create a struct to thold a comment
	// we use a struct tag [``] to maek the names display in lowercase
	var incomindData struct {
		Content string `json:"content"`
		Author  string `json:"author"`
	}

	// Perform the decoding
	err := json.NewDecoder(r.Body).Decode(&incomindData)
	if err != nil {
		a.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Fprintf(w, "%+v\n", incomindData)
}
