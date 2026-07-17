package handlers

import (
	"fmt"
	"log"
	"net/http"
)

// errorResponse writes a JSON error response with the given HTTP status code.
func (h *APIHandler) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	msg := map[string]any{"error": message}

	if err := h.writeJSON(w, status, msg, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// serverErrorResponse logs the error and returns a generic 500 Internal Server Error response.
func (h *APIHandler) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("server error: %v", err)
	message := "the server encountered a problem and could not process your request"
	h.errorResponse(w, r, http.StatusInternalServerError, message)
}

// notFoundResponse returns a 404 Not Found error response.
func (h *APIHandler) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	h.errorResponse(w, r, http.StatusNotFound, message)
}

// methodNotAllowedResponse returns a 405 Method Not Allowed error response.
func (h *APIHandler) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	h.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// badRequestResponse returns a 400 Bad Request error response.
func (h *APIHandler) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// failedValidationResponse returns a 422 Unprocessable Entity response with validation errors.
func (h *APIHandler) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	h.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
