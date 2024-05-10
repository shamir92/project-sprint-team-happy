package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
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

	if statusCode == http.StatusInternalServerError {
		errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

		errLog.Output(2, trace)
	}

	s.errorResponse(w, r, statusCode, err)
}
