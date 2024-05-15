package main

import (
	"halosuster/configuration"
	"halosuster/internal/database"
	"halosuster/internal/helper"
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
	var bcryptConfiguration configuration.IBcryptConfiguration = configuration.NewBcryptConfiguration()

	app := fiber.New()

	// For External Interfaces
	postgresWriter, err := database.NewPostgresWriter(dbConfiguration)
	if err != nil {
		log.Fatal(err)
	}

	bcrypt := helper.NewBcryptPasswordHash(bcryptConfiguration)
	jwtManager := helper.NewJwt(jwtConfiguration)

	// Routes
	publicRouteParam := route.PublicRouteParams{
		App:              app,
		AppConfiguration: appConfiguration,
		PostgresWriter:   postgresWriter,
		JwtConfiguration: jwtConfiguration,
		HelperBcrypt:     bcrypt,
		JWTManager:       jwtManager,
	}
	route.PublicRoutes(publicRouteParam)
	route.PrivateRoutes(route.PrivateRoutesParam{
		App:              app,
		AppConfiguration: appConfiguration,
		PostgresWriter:   postgresWriter,
		JwtManager:       jwtManager,
	})
	//nolint:errcheck
	log.Fatal(app.Listen(appConfiguration.GetPort()))
}
