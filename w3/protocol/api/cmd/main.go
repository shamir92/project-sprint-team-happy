package main

import (
	"halosuster/configuration"
	"halosuster/protocol/api/route"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	// For configuration
	var appConfiguration configuration.IAppConfiguration = configuration.NewAppConfiguration()

	var dbConfiguration configuration.IDatabaseWriter = configuration.NewDatabaseWriter()
	var jwtConfiguration configuration.IJWTConfiguration = configuration.NewJWTConfiguration()

	app := fiber.New()

	// For External Interfaces
	// var postgresWriter postgres.IPostgresWriter = database.NewPostgresWriter(dbConfiguration)

	// Routes
	publicRouteParam := route.PublicRouteParams{
		App:                   app,
		AppConfiguration:      appConfiguration,
		DatabaseConfiguration: dbConfiguration,
		JwtConfiguration:      jwtConfiguration,
	}
	route.PublicRoutes(publicRouteParam)
	//nolint:errcheck
	log.Fatal(app.Listen(appConfiguration.GetPort()))
}
