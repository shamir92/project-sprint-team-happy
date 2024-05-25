package controller

import (
	"belimang/domain/usecase"
	"belimang/internal/helper"
	"belimang/protocol/api/dto"

	"github.com/gofiber/fiber/v2"
)

type orderController struct {
	orderUsecase usecase.IOrderUsecase
}

func NewOrderController(orderUsecase usecase.IOrderUsecase) *orderController {
	return &orderController{orderUsecase}
}

func (oc *orderController) PostOrderEstimate(ctx *fiber.Ctx) error {
	var body usecase.MakeOrderEstimatePayload

	if err := ctx.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	if err := helper.ValidateStruct(body); err != nil {
		return err
	}

	user := ctx.Locals("user").(*helper.JsonWebTokenClaims)

	order, err := oc.orderUsecase.MakeOrderEstimate(body, user.UserID)

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": dto.OrderEstimateResponseDto{
			OrderId:                        order.ID,
			TotalPrice:                     order.TotalPrice,
			EstimatedDeliveryTimeInMinutes: order.EstimatedDeliveryTime,
		},
	})
}
