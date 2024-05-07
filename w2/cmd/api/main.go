package main

import (
	"eniqlostore/internal/postgres"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	db, err := postgres.NewDB("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	if err != nil {
		log.Fatalf("database failed to open: %v", err)
	}

	server := newServer(serverOpts{
		db: db,
	})

	httpServer := server.newHttpServer(":8080")

	err = httpServer.ListenAndServe()

	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
