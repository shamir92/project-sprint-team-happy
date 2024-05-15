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

type PrivateRouteParams struct {
	App              *fiber.App
	AppConfiguration configuration.IAppConfiguration
	PostgresWriter   database.IPostgresWriter
	JwtConfiguration configuration.IJWTConfiguration
	HelperBcrypt     helper.IBcryptPasswordHash
	JWTManager       helper.IJWTManager
	S3Writer         database.IS3Writer
}

// TODO : add routes to here.
func PrivateRoutes(params PrivateRouteParams) {
	// TODO: initiation of repository
	// var userRepository repository.IUserRepository = repository.NewUserRepository(params.PostgresWriter.GetDB())
	var s3Repository repository.IS3Repository = repository.NewS3Repository(params.S3Writer)

	// TODO: initiation of usecase/ service
	var pingUsecase usecase.IPingUsecase = usecase.NewPingUsecase()
	// var userITUsecase usecase.IUserITUsecase = usecase.NewUserITUsecase(params.HelperBcrypt, userRepository, params.JWTManager)
	var s3Usecase usecase.IImageUsecase = usecase.NewImageUsecase(s3Repository)

	// TODO: initiation of controller/ handler
	var pingController controller.IPingController = controller.NewPingController(pingUsecase)
	// var userITController controller.IUserITController = controller.NewUserITController(userITUsecase)
	var imageController controller.IImageController = controller.NewImageController(s3Usecase)

	// Create routes group.
	route := params.App.Group("/v1")
	route.Get("/ping", pingController.GetPingController)
	route.Post("/image", imageController.UploadImage)
	// route.Post("/user/it/login", userITController.LoginUserIT)

	//
	log.Println(route)

}
