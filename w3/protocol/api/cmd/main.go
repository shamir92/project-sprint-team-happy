package main

import (
	"halosuster/configuration"
	"halosuster/protocol/api/route"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {

	// For configuration
	var appConfiguration configuration.IAppConfiguration = configuration.NewAppConfiguration()

	// var dbConfiguration configuration.IDatabaseWriter = configuration.NewDatabaseWriter()
	// var jwtConfiguration configuration.IJWTConfiguration = configuration.NewJWTConfiguration()

	app := fiber.New()

	// For External Interfaces
	// var postgresWriter postgres.IPostgresWriter = database.NewPostgresWriter(dbConfiguration)

	// Routes
	route.PublicRoutes(app)
	//nolint:errcheck
	log.Fatal(app.Listen(appConfiguration.GetPort()))
}
