package main

import (
	"eniqlostore/internal/postgres"
	"eniqlostore/internal/repository"
	"eniqlostore/internal/service"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {

	db, err := postgres.NewDB("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	if err != nil {
		log.Fatalf("database failed to open: %v", err)
	}

	service := service.NewService(service.ServiceDeps{
		UserRepository: repository.NewUserRepository(db),
	})

	httpServer := newServer(service)

	err = http.ListenAndServe(":8080", httpServer.router)

	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
