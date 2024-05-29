package controller

import (
	"belimang/domain/usecase"
	"belimang/internal/helper"
	"belimang/protocol/api/dto"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
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
	context := ctx.Locals("ctx").(context.Context)
	tracer := ctx.Locals("tracer").(trace.Tracer)

	_, span := tracer.Start(context, "CreateItem")
	defer span.End()
	var body usecase.CreateMerchanItemPayload

	if err := ctx.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	if err := helper.ValidateStruct(body); err != nil {
		return err
	}

	body.MerchantID = ctx.Params("merchantId")

	user := ctx.Locals("user").(*helper.JsonWebTokenClaims)

	item, err := mic.itemUsecase.Create(context, body, user.UserID)

	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"itemId": item.ID,
		},
	})
}

func (mic *merchantItemController) GetItems(ctx *fiber.Ctx) error {
	context := ctx.Locals("ctx").(context.Context)
	tracer := ctx.Locals("tracer").(trace.Tracer)

	_, span := tracer.Start(context, "GetItems")
	defer span.End()
	var query dto.FindMerchantItemPayload = dto.FindMerchantItemPayload{
		MerchantID:  ctx.Params("merchantId"),
		Limit:       ctx.Query("limit"),
		Offset:      ctx.Query("offset"),
		Name:        ctx.Query("name"),
		ItemID:      ctx.Query("itemId"),
		SortCreated: ctx.Query("createdAt"),
		Category:    ctx.Query("category"),
	}

	rows, meta, err := mic.itemUsecase.FindItemsByMerchant(context, query)

	if err != nil {
		return err
	}

	var items []dto.MerchanItemDto = make([]dto.MerchanItemDto, 0)

	for _, row := range rows {
		items = append(items, dto.MerchanItemDto{
			ID:        row.ID,
			Name:      row.Name,
			Category:  row.Category.String(),
			CreatedAt: row.CreatedAt,
			Price:     row.Price,
			ImageUrl:  row.ImageUrl,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": items,
		"meta": meta,
	})
}
