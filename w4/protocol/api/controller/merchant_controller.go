package controller

import (
	"belimang/domain/usecase"
	"belimang/internal/helper"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type merchantController struct {
	merchantUsecase usecase.IMerchantUsecase
}

func NewMerchantController(merchantUsecase usecase.IMerchantUsecase) *merchantController {
	return &merchantController{merchantUsecase}
}

func (c *merchantController) Browse(ctx *fiber.Ctx) error {
	var query usecase.MerchantFetchQuery

	if err := ctx.QueryParser(&query); err != nil {
		return fiber.ErrBadRequest
	}

	response, err := c.merchantUsecase.Fetch(query)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    response,
	})
}

func (c *merchantController) Create(ctx *fiber.Ctx) error {
	var payload usecase.MerchantCreatePayload

	if err := ctx.BodyParser(&payload); err != nil {
		return fiber.ErrBadRequest
	}

	if err := helper.ValidateStruct(payload); err != nil {
		return err
	}

	response, err := c.merchantUsecase.Create(payload)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"data":    response,
	})
}
