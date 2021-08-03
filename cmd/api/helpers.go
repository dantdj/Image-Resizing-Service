package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// A wrapper for an object to be returned as JSON in a response
type envelope map[string]interface{}

// Takes the destination http.ResponseWriter, the HTTP status code to send,
// the data to encode to JSON, and a header map containing any additional
// HTTP headers to include in the response, and writes the JSON object
// to a given ResponseWriter
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	// At this point, we know that we won't encounter any more errors before writing the
	// response, so it's safe to add any headers that we want to include.
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// Reads an int value from the query string of a request. If the value is missing or an error occurs,
// it will return the provided default value.
func (app *application) readInt(queryStr url.Values, key string, defaultValue int) int {
	str := queryStr.Get(key)

	if str == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(str)
	if err != nil {
		log.WithFields(log.Fields{
			"stringValue": str,
			"error":       err,
			"queryString": queryStr,
		}).Error("failed extracting integer value from query string")
		return defaultValue
	}

	return value
}
