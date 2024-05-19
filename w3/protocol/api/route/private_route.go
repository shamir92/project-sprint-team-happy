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
	var userRepository repository.IUserRepository = repository.NewUserRepository(params.PostgresWriter.GetDB())
	var medicalRecordPatientRepository repository.IMedicalRecordPatientRepository = repository.NewMedicalRecordPatientRepository(params.PostgresWriter.GetDB())
	var s3Repository repository.IS3Repository = repository.NewS3Repository(params.S3Writer)
	var medicalRecordRepo repository.IMedicalRecordRepository = repository.NewMedicalRecordRepository(params.PostgresWriter.GetDB())

	// TODO: initiation of usecase/ service
	var pingUsecase usecase.IPingUsecase = usecase.NewPingUsecase()
	var userITUsecase usecase.IUserITUsecase = usecase.NewUserITUsecase(params.HelperBcrypt, userRepository, params.JWTManager)
	var s3Usecase usecase.IImageUsecase = usecase.NewImageUsecase(s3Repository)
	var nurseUseCase = usecase.NewUserNurseUseCase(userRepository, params.HelperBcrypt, params.JWTManager)
	var medicalRecordPatientUsecase usecase.IMedicalRecordPatientUsecase = usecase.NewMedicalRecordPatientUsecase(medicalRecordPatientRepository)
	var medicalRecordUsecase usecase.IMedicalRecordUsecase = usecase.NewMedicalRecordUsecase(medicalRecordPatientRepository, medicalRecordRepo)

	// TODO: initiation of controller/ handler
	var pingController controller.IPingController = controller.NewPingController(pingUsecase)
	var userITController controller.IUserITController = controller.NewUserITController(userITUsecase)
	var imageController controller.IImageController = controller.NewImageController(s3Usecase)
	var nurseController = controller.NewUserNurseController(nurseUseCase)
	var medicalRecordPatientController controller.IMedicalRecordPatientController = controller.NewMedicalRecordPatientController(medicalRecordPatientUsecase)
	var medicalRecordController controller.IMedicalRecordController = controller.NewMedicalRecordController(medicalRecordUsecase)

	// Create routes group.
	route := params.App.Group("/v1")
	route.Use(middleware.AuthMiddleware(params.JWTManager))
	route.Get("/ping", pingController.GetPingController)
	route.Post("/image", imageController.UploadImage)
	route.Get("/user", userITController.GetListUsers)
	route.Route("/user/nurse", func(router fiber.Router) {
		router.Post("/register", nurseController.CreateNurse)
		router.Put("/:userNurseId", nurseController.UpdateNurse)
		router.Delete("/:userNurseId", nurseController.DeleteNurse)
		router.Post("/:userNurseId/access", nurseController.SetAccessNurse)
	})

	medical := route.Group("/medical")
	medical.Post("/patient", medicalRecordPatientController.Create)

	medical.Get("/patient", medicalRecordPatientController.Browse)

	medical.Post("/record", medicalRecordController.Create)
	medical.Get("/record", medicalRecordController.GetMedicalRecords)
}
