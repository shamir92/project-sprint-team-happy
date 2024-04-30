package main

import (
	"fmt"
	"gin-mvc/internal"
	"gin-mvc/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	// loading .env to golang
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// loading db postgres
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbParams := os.Getenv("DB_PARAMS")
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s %s", dbUsername, dbPassword, dbHost, dbPort, dbName, dbParams)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	fmt.Println("Successfully connected!")
	// load the memory to global package
	internal.SetDB(db)

	// loading gin to go
	r := gin.Default()
	routes.RegisterRoutes(r)
	r.Run() // listen and serve on 0.0.0.0:8080
}
