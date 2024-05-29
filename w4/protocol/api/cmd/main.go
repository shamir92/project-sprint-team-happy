package main

import (
	"belimang/configuration"
	"belimang/internal/database"
	"belimang/internal/helper"
	"belimang/protocol/api/controller"
	"belimang/protocol/api/route"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("error loading .env file")
	}

	// For configuration
	var jaegerConfiguration configuration.IJaegerConfiguration = configuration.NewJaegerConfiguration()
	var appConfiguration configuration.IAppConfiguration = configuration.NewAppConfiguration()
	var dbConfiguration configuration.IDatabaseWriter = configuration.NewDatabaseWriter()
	var jwtConfiguration configuration.IJWTConfiguration = configuration.NewJWTConfiguration()
	var bcryptConfiguration configuration.IBcryptConfiguration = configuration.NewBcryptConfiguration()

	app := fiber.New(fiber.Config{
		ErrorHandler: controller.ErrorHandler,
	})
	// app.Use(otelfiber.Middleware())
	tracer := jaegerConfiguration.GetTracer()
	app.Use(func(c *fiber.Ctx) error {
		ctx, span := tracer.Start(context.Background(), c.Path())
		defer span.End()

		c.Locals("ctx", ctx)
		c.Locals("tracer", tracer)
		return c.Next()
	})

	app.Use(logger.New())
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

	privateRouteParam := route.PrivateRouteParam{
		App:              app,
		AppConfiguration: appConfiguration,
		PostgresWriter:   postgresWriter,
		JwtConfiguration: jwtConfiguration,
		HelperBcrypt:     bcrypt,
		JWTManager:       jwtManager,
	}
	route.PrivateRoutes(privateRouteParam)

	//nolint:errcheck
	log.Fatal(app.Listen(appConfiguration.GetPort()))
}
