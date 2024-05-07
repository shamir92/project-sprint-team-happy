package httpserver

import (
	"net/http"
)

type httpStatusCodeProvider interface {
	HTTPStatusCode() int
}

func (s *HttpServer) errorResponse(w http.ResponseWriter, r *http.Request, status int, err error) {
	statusCode := http.StatusInternalServerError

	switch e := err.(type) {
	case httpStatusCodeProvider:
		statusCode = e.HTTPStatusCode()
	}

	s.writeJSON(w, r, statusCode, map[string]any{"error": err.Error()})
}
