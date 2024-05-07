package main

import (
	"net/http"
)

func (s *server) errorResponse(w http.ResponseWriter, r *http.Request, status int, err error) {
	env := map[string]any{"error": err.Error()}

	s.writeJSON(w, r, status, env)
}
