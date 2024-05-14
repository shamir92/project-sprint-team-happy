package route

import (
	"halosuster/configuration"
	"halosuster/protocol/api/controller"
	"log"

	"github.com/gofiber/fiber/v2"
)

type PublicRouteParams struct {
	App                   *fiber.App
	AppConfiguration      configuration.IAppConfiguration
	DatabaseConfiguration configuration.IDatabaseWriter
	JwtConfiguration      configuration.IJWTConfiguration
}

// TODO : add routes to here.
func PublicRoutes(params PublicRouteParams) {
	log.Println(params.AppConfiguration)
	log.Println(params.DatabaseConfiguration)
	log.Println(params.JwtConfiguration)
	// TODO: initiation of repository

	// TODO: initiation of usecase/ service

	// TODO: initiation of controller/ handler
	var pingController controller.IPingController = controller.NewPingController()

	// Create routes group.
	route := params.App.Group("/api/v1")
	route.Get("/ping", pingController.GetPingController)

	//
	log.Println(route)

}
