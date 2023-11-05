package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	// Convert the notFoundResponse() helper to a http.Handler using the
	// http.HandlerFunc() adapter, and then set it as the custom error handler for 404
	// Not Found responses.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	// Likewise, convert the methodNotAllowedResponse() helper to a http.Handler and set
	// it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/edtoys", app.listEdToysHandler)
	router.HandlerFunc(http.MethodPost, "/v1/edtoys", app.createEdtoysHandler)
	router.HandlerFunc(http.MethodGet, "/v1/edtoys/:id", app.showEdtoysHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/edtoys/:id", app.updateEdToysHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/edtoys/:id", app.deleteEdToysHandler)
	return app.recoverPanic(app.rateLimit(router))
}
