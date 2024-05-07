package main

import (
	"eniqlostore/internal/postgres"
	httpserver "eniqlostore/server/http"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	db, err := postgres.NewDB("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	if err != nil {
		log.Fatalf("database failed to open: %v", err)
	}

	httpServer := httpserver.New(httpserver.ServerOpts{
		DB:   db,
		Addr: ":8080",
	})

	server := httpServer.Server()

	err = server.ListenAndServe()

	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
