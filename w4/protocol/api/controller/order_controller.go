package controller

import (
	"belimang/domain/usecase"
	"belimang/internal/helper"
	"belimang/protocol/api/dto"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

type orderController struct {
	orderUsecase usecase.IOrderUsecase
}

func NewOrderController(orderUsecase usecase.IOrderUsecase) *orderController {
	return &orderController{orderUsecase}
}

func (oc *orderController) PostOrderEstimate(ctx *fiber.Ctx) error {
	context := ctx.Locals("ctx").(context.Context)
	tracer := ctx.Locals("tracer").(trace.Tracer)

	_, span := tracer.Start(context, "PostOrderEstimate")
	defer span.End()
	var body usecase.MakeOrderEstimatePayload

	if err := ctx.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	if err := helper.ValidateStruct(body); err != nil {
		return err
	}

	user := ctx.Locals("user").(*helper.JsonWebTokenClaims)

	order, err := oc.orderUsecase.MakeOrderEstimate(context, body, user.UserID)

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

func (oc *orderController) PlaceOrder(ctx *fiber.Ctx) error {
	context := ctx.Locals("ctx").(context.Context)
	tracer := ctx.Locals("tracer").(trace.Tracer)

	_, span := tracer.Start(context, "PlaceOrder")
	defer span.End()
	var body dto.PlaceOrderRequestDto

	if err := ctx.BodyParser(&body); err != nil {
		return fiber.ErrBadGateway
	}

	user := ctx.Locals("user").(*helper.JsonWebTokenClaims)

	order, err := oc.orderUsecase.PlaceOrder(context, body.OrderId, user.UserID)

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": dto.PlaceOrderResponseDto{
			OrderId: order.ID.String(),
		},
	})
}

func (oc *orderController) GetUserOrders(ctx *fiber.Ctx) error {
	context := ctx.Locals("ctx").(context.Context)
	tracer := ctx.Locals("tracer").(trace.Tracer)

	_, span := tracer.Start(context, "GetUserOrders")
	defer span.End()
	var body dto.GetOrderSearchParams

	if err := ctx.QueryParser(&body); err != nil {
		log.Printf("ERROR | GetUserOrders() | %v\n", err)
		return err
	}

	user := ctx.Locals("user").(*helper.JsonWebTokenClaims)

	orders, err := oc.orderUsecase.GetOrders(context, body, user.UserID)

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": orders,
	})
}
