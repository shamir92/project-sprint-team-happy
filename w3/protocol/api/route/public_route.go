package route

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")
	log.Println(route)
	// Routes for GET method:
	// route.Get("/short-url/:shortUrl", controllers.GetShortUrl) // get list of all books

	// Routes for POST method:
	// route.Post("/short-url", controllers.CreateShortUrl) // get list of all books

	// route.Post("/user/sign/in", controllers.UserSignIn) // auth, return Access & Refresh tokens
}
