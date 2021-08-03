package main

import (
	"fmt"
	"net/http"
)

// Recovers from a panicked request, returning a `serverErrorResponse()` to the user so that the response
// isn't empty when a panic occurs. This middleware will only affect the same goroutine that executed it.
// Further goroutines created by handlers will NOT be covered, and must handle panics themselves.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
