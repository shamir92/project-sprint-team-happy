package controller

import (
	"belimang/domain/usecase"
	"belimang/internal/helper"
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

type merchantController struct {
	merchantUsecase usecase.IMerchantUsecase
}

func NewMerchantController(merchantUsecase usecase.IMerchantUsecase) *merchantController {
	return &merchantController{merchantUsecase}
}

func (c *merchantController) Browse(ctx *fiber.Ctx) error {
	context := ctx.Locals("ctx").(context.Context)
	tracer := ctx.Locals("tracer").(trace.Tracer)

	_, span := tracer.Start(context, "Browse")
	defer span.End()
	var query usecase.MerchantFetchQuery

	if err := ctx.QueryParser(&query); err != nil {
		return fiber.ErrBadRequest
	}

	response, err := c.merchantUsecase.Fetch(context, query)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    response,
	})
}

func (c *merchantController) Create(ctx *fiber.Ctx) error {
	context := ctx.Locals("ctx").(context.Context)
	tracer := ctx.Locals("tracer").(trace.Tracer)

	_, span := tracer.Start(context, "Create")
	defer span.End()
	var payload usecase.MerchantCreatePayload

	if err := ctx.BodyParser(&payload); err != nil {
		return fiber.ErrBadRequest
	}

	if err := helper.ValidateStruct(payload); err != nil {
		return err
	}

	response, err := c.merchantUsecase.Create(context, payload)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"data":    response,
	})
}
