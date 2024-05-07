package main

import (
	"eniqlostore/internal/postgres"
	httpserver "eniqlostore/server/http"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func connectionString() string {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbParams := os.Getenv("DB_PARAMS")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", dbUsername, dbPassword, dbHost, dbPort, dbName, dbParams)

	return connStr
}

func main() {
	godotenv.Load()

	db, err := postgres.NewDB(connectionString())

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
