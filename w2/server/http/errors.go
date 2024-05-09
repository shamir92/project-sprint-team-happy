package httpserver

import (
	"net/http"
)

type httpStatusCodeProvider interface {
	HTTPStatusCode() int
}

// TODO: Re-thinking how to structure this function
func (s *HttpServer) errorResponse(w http.ResponseWriter, r *http.Request, status int, err error) {
	s.writeJSON(w, r, status, map[string]any{"error": err.Error()})
}

func (s *HttpServer) errorBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	s.writeJSON(w, r, http.StatusBadRequest, map[string]any{"error": err.Error()})
}

func (s *HttpServer) handleError(w http.ResponseWriter, r *http.Request, err error) {
	var statusCode int

	switch e := err.(type) {
	case httpStatusCodeProvider:
		statusCode = e.HTTPStatusCode()
	default:
		statusCode = http.StatusInternalServerError
	}

	s.errorResponse(w, r, statusCode, err)
}
