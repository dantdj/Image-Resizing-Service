package main

import (
	"fmt"
	"net/http"
)

// Generic helper for logging an error message.
func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

// The errorResponse() method is a generic helper for sending JSON-formatted error
// messages to the client with a given status code
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	// Write the response using the writeJSON() helper. If this returns an
	// error then log it, and fall back to sending the client an empty response with a
	// 500 Internal Server Error status code.
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// Logs the detailed error message, then uses the errorResponse() helper to send
// a 500 Internal Server Error status code and JSON response (containing a generic
// error message) to the client.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// Used to send a 404 Not Found status code and JSON response to the client.
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// Used to send a 405 Method Not Allowed status code and JSON response to the client.
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// Used to send a 400 Bad Request status code and JSON response to the client.
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// Used to send a 415 Unsupported Media Type status code and JSON response to the client
func (app *application) unsupportedMediaTypeResponse(w http.ResponseWriter, r *http.Request) {
	message := "unsupported media type given"
	app.errorResponse(w, r, http.StatusUnsupportedMediaType, message)
}
