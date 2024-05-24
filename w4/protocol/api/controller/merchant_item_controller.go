package controller

import (
	"belimang/domain/usecase"
	"belimang/internal/helper"

	"github.com/gofiber/fiber/v2"
)

type merchantItemController struct {
	itemUsecase usecase.IMerchantItemUsecase
}

func NewMerchantItemController(itemUsecase usecase.IMerchantItemUsecase) *merchantItemController {
	return &merchantItemController{
		itemUsecase: itemUsecase,
	}
}

func (mic *merchantItemController) CreateItem(ctx *fiber.Ctx) error {
	var body usecase.CreateMerchanItemPayload

	if err := ctx.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	if err := helper.ValidateStruct(body); err != nil {
		return err
	}

	body.MerchantID = ctx.Params("merchantId")

	user := ctx.Locals("user").(*helper.JsonWebTokenClaims)

	item, err := mic.itemUsecase.Create(body, user.UserID)

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"itemId": item.ID,
		},
	})
}
