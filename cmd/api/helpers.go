package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func (app *application) versioncontrolURI(pattern string) string {
	return fmt.Sprintf(`/v%v/%v`, app.config.version, pattern)
}

func (a *application) readJson(w http.ResponseWriter, r *http.Request, dest any) error {

	// what is the max size of the request body (250KB seems reasonable)
	maxBytes := 256_000
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	// our decoder will check for unknown fields
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	// let start the decoding
	err := dec.Decode(dest)
	if err != nil {
		// Check for different errors
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("the body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("the body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("the body contains the incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("the body contains the incorrect  JSON type (at character %d)", unmarshalTypeError.Offset)
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
	// almost done. Let's lastly check if there is any data after
	// the valid JSON data. Maybe the person is trying to send
	// multiple request bodies during one request
	// We call decode once more to see if it gives us back anything
	// we use a throw away struct 'struct{}{}' to hold the result
	err = dec.Decode(&struct{}{})

	if !errors.Is(err, io.EOF) {
		// there is more data present
		return errors.New("the body must only contain a single JSON value")
	}

	return nil
}
