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

type PrivateRouteParam struct {
	App              *fiber.App
	AppConfiguration configuration.IAppConfiguration
	PostgresWriter   database.IPostgresWriter
	JwtConfiguration configuration.IJWTConfiguration
	HelperBcrypt     helper.IBcryptPasswordHash
	JWTManager       helper.IJWTManager
}

func PrivateRoutes(params PrivateRouteParam) {
	var merchantRepository repository.IMerchantRepository = repository.NewMerchantRepository(params.PostgresWriter.GetDB())
	var merchantUsecase usecase.IMerchantUsecase = usecase.NewMerchantUsecase(merchantRepository)
	var merchantController = controller.NewMerchantController(merchantUsecase)

	var merchantItemRepository repository.IMerchantItemRepository = repository.NewMerchanItemRepository(params.PostgresWriter.GetDB())
	var merchantItemUsecase usecase.IMerchantItemUsecase = usecase.NewMerchanItemUsecase(merchantItemRepository)
	var merchantItemController = controller.NewMerchantItemController(merchantItemUsecase)

	route := params.App.Group("v1/admin")

	merchant := route.Group("merchants")
	merchant.Get("/", merchantController.Browse)
	merchant.Post("/", merchantController.Create)
	merchant.Post("/:merchantId/items", merchantItemController.CreateItem)
	merchant.Get("/:merchantId/items", merchantItemController.CreateItem)

}
