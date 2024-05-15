package route

import (
	"halosuster/configuration"
	"halosuster/domain/repository"
	"halosuster/domain/usecase"
	"halosuster/internal/database"
	"halosuster/internal/helper"
	"halosuster/protocol/api/controller"
	"log"

	"github.com/gofiber/fiber/v2"
)

type PublicRouteParams struct {
	App              *fiber.App
	AppConfiguration configuration.IAppConfiguration
	PostgresWriter   database.IPostgresWriter
	JwtConfiguration configuration.IJWTConfiguration
	HelperBcrypt     helper.IBcryptPasswordHash
	JWTManager       helper.IJWTManager
}

// TODO : add routes to here.
func PublicRoutes(params PublicRouteParams) {
	log.Println(params.AppConfiguration)
	log.Println(params.PostgresWriter)
	log.Println(params.JwtConfiguration)
	// TODO: initiation of repository
	var userRepository repository.IUserRepository = repository.NewUserRepository(params.PostgresWriter.GetDB())

	// TODO: initiation of usecase/ service
	var pingUsecase usecase.IPingUsecase = usecase.NewPingUsecase()
	var userITUsecase usecase.IUserITUsecase = usecase.NewUserITUsecase(params.HelperBcrypt, userRepository, params.JWTManager)
	// TODO: initiation of controller/ handler
	var pingController controller.IPingController = controller.NewPingController(pingUsecase)
	var userITController controller.IUserITController = controller.NewUserITController(userITUsecase)

	// Create routes group.
	route := params.App.Group("/v1")
	route.Get("/ping", pingController.GetPingController)
	route.Post("/user/it/register", userITController.RegisterUserIT)

	//
	log.Println(route)

}
