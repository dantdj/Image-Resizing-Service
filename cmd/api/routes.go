package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/ping", app.pingHandler)
	router.HandlerFunc(http.MethodGet, "/capabilities", app.capabilitiesHandler)

	router.HandlerFunc(http.MethodPost, "/resize", app.resizeHandler)

	return app.recoverPanic(router)
}
