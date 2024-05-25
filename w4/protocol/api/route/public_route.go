package route

import (
	"belimang/configuration"
	"belimang/domain/repository"
	"belimang/domain/usecase"
	"belimang/internal/database"
	"belimang/internal/helper"
	"belimang/protocol/api/controller"

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

func PublicRoutes(params PublicRouteParams) {

	var userRepository repository.IUserRepository = repository.NewUserRepository(params.PostgresWriter.GetDB())
	var userUsecase usecase.IUserUsecase = usecase.NewUserUsecase(userRepository, params.JWTManager, params.HelperBcrypt)
	var userController = controller.NewUserController(userUsecase)

	v1 := params.App.Group("v1")

	v1.Route("/users", func(router fiber.Router) {
		router.Post("/register", userController.Register)
	})
}