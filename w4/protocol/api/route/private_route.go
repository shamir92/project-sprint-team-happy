package route

import (
	"belimang/configuration"
	"belimang/domain/repository"
	"belimang/domain/usecase"
	"belimang/internal/database"
	"belimang/internal/helper"
	"belimang/protocol/api/controller"
	"belimang/protocol/api/middleware"

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

	var orderRepository repository.IOrderRepository = repository.NewOrderRepository(params.PostgresWriter.GetDB())
	var orderUsecase usecase.IOrderUsecase = usecase.NewOrderUsecase(orderRepository, merchantItemRepository)
	var orderController = controller.NewOrderController(orderUsecase)

	params.App.Use(middleware.AuthMiddleware(params.JWTManager))
	params.App.Route("/admin/merchants", func(merchant fiber.Router) {
		merchant.Get("/", merchantController.Browse)
		merchant.Post("/", merchantController.Create)
		merchant.Post("/:merchantId/items", merchantItemController.CreateItem)
		merchant.Get("/:merchantId/items", merchantItemController.GetItems)
	})

	params.App.Route("/users", func(router fiber.Router) {
		router.Post("/estimate", orderController.PostOrderEstimate)
		router.Post("/orders", orderController.PlaceOrder)
		router.Get("/orders", orderController.GetUserOrders)
	})

	params.App.Get("/merchants/nearby/:latlon", merchantController.BrowseNearby)
	params.App.Get("/merchants/nearby/*", merchantController.BrowseNearbyInvalid)
}
