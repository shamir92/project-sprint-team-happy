package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *server) routes() {
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hey, What's Up!"))
	})

	s.router.Route("/v1", func(r chi.Router) {
		r.Get("/staff/register", s.handleCreateStaff)
	})
}
