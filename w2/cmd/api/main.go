package main

import (
	"eniqlostore/internal/postgres"
	"eniqlostore/internal/service"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	db := postgres.NewDB("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	service := service.NewService(db)
	httpServer := newServer(service)

	http.ListenAndServe(":8080", httpServer.router)
}
