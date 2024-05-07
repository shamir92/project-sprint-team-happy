package main

import (
	"eniqlostore/internal/service"

	"github.com/go-chi/chi/v5"
)

type server struct {
	service *service.Service
	router  *chi.Mux
}

func newServer(service *service.Service) *server {
	s := &server{
		service: service,
	}

	s.router = chi.NewRouter()

	s.routes()

	return s
}
