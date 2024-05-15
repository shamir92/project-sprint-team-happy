package route

import (
	"halosuster/configuration"
	"halosuster/domain/repository"
	"halosuster/domain/usecase"
	"halosuster/internal/database"
	"halosuster/internal/helper"
	"halosuster/protocol/api/controller"
	"halosuster/protocol/api/middleware"

	"github.com/gofiber/fiber/v2"
)

type PrivateRoutesParam struct {
	App              *fiber.App
	AppConfiguration configuration.IAppConfiguration
	PostgresWriter   database.IPostgresWriter
	JwtManager       helper.IJWTManager
}

func PrivateRoutes(params PrivateRoutesParam) {
	var userRepository repository.IUserRepository = repository.NewUserRepository(params.PostgresWriter.GetDB())

	var nurseUseCase = usecase.NewUserNurseUseCase(userRepository)
	// TODO: initiation of controller/ handler
	var nurseController = controller.NewUserNurseController(nurseUseCase)

	privateV1Router := params.App.Group("v1")

	privateV1Router.Use(middleware.AuthMiddleware(params.JwtManager))
	privateV1Router.Route("/user/nurse", func(router fiber.Router) {
		router.Post("/register", nurseController.CreateNurse)
	})
}
