package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Errorw("internal server error", "method", r.Method, "url", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Warnf("Bad request error", "method", r.Method, "url", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) conflictErr(w http.ResponseWriter, r *http.Request, err error) {
	
	app.logger.Errorf("Conflict error:", "method", r.Method, "url", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Warnf("Request Not Found", "method", r.Method, "url", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusNotFound, "request not found")
}

func (app *application) unAuthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("Unauthorized error", "method", r.Method, "url", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) BasicunAuthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("Unauthorized error", "method", r.Method, "url", r.URL.Path, "error", err.Error())
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}


func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw("Forbidden", "method", r.Method, "url", r.URL.Path)

	writeJSONError(w, http.StatusForbidden, "Forbidden")
}


func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}