// file: cmd/api/helpers.go
// Description: Helper functions for JSON encoding/decoding and error handling

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// creating an envelope type
type envelope map[string]any

// Helper function to write json, has the following parameters:
// response writer, status code for the server, data of custom type envelope which is a map to encode, and the headers to specify.
// returns an error
func (a *appDependencies) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {

	// encodes data into json format by using indenting for better readability
	jsResponse, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// append json and add a new line after each appendage
	jsResponse = append(jsResponse, '\n')

	// add any headers that we want to the response
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

// Helper function to read json from request body also performs error handling
// takes in response writer, request and a destination of any type
// returns an error
func (a *appDependencies) readJson(w http.ResponseWriter, r *http.Request, dest any) error {

	// limit the size of the request body to 256000 bytes
	maxBytes := 256_000
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// our decoder will check for unknown fields
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// let start the decoding
	err := dec.Decode(dest)

	// Check for different errors
	if err != nil {
		// syntax error
		var syntaxError *json.SyntaxError
		// incorrect type error
		var unmarshalTypeError *json.UnmarshalTypeError
		// empty body error
		var invalidUnmarshalError *json.InvalidUnmarshalError
		// max size error
		var maxBytesError *http.MaxBytesError

		// using a switch to handle different errors
		switch {
		// check for syntax error
		case errors.As(err, &syntaxError):
			return fmt.Errorf("the body contains badly-formed JSON (at character %d)", syntaxError.Offset)
			// check for unexpected EOF error
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("the body contains badly-formed JSON")
			// check for incorrect type error
		case errors.As(err, &unmarshalTypeError):
			// if the field is not empty, it means we have a specific field that is incorrect
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("the body contains the incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("the body contains the incorrect  JSON type (at character %d)", unmarshalTypeError.Offset)
			// check for empty body error
		case errors.Is(err, io.EOF):
			return errors.New("the body must not be empty")
			// check for unknown field error
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(),
				"json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
			//Size
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("the body must not be larger than %d bytes", maxBytesError.Limit)
			// some error the programmer made
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	// call decode again to check if there is only a single json value in the body
	err = dec.Decode(&struct{}{})

	// if the error is not EOF, then there is more than one value in the body
	if !errors.Is(err, io.EOF) {
		return errors.New("the body must only contain a single JSON value")
	}

	return nil
}

// Helper function to read an id parameter from the url
func (a *appDependencies) readIDParam(r *http.Request) (int64, error) {
	// get the url parameters
	params := httprouter.ParamsFromContext(r.Context())

	// convert the id parameter to an int64
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}
