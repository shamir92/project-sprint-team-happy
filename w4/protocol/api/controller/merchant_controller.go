package controller

import (
	"belimang/domain/usecase"
	"belimang/internal/helper"
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

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

	response, paginationMeta, err := c.merchantUsecase.Fetch(context, query)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    response,
		"meta":    paginationMeta,
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

	return ctx.Status(http.StatusCreated).JSON(response)
}

func (c *merchantController) BrowseNearby(ctx *fiber.Ctx) error {
	context := ctx.Locals("ctx").(context.Context)
	tracer := ctx.Locals("tracer").(trace.Tracer)

	_, span := tracer.Start(context, "Browse")
	defer span.End()
	coordinate := strings.Split(ctx.Params("latlon"), ",")
	coordinate[0] = strings.TrimSpace(coordinate[0])
	coordinate[1] = strings.TrimSpace(coordinate[1])

	lat, err := strconv.ParseFloat(coordinate[0], 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	lon, err := strconv.ParseFloat(coordinate[1], 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	log.Println(lat)
	log.Println(lon)
	var query usecase.MerchantFetchQuery

	if err := ctx.QueryParser(&query); err != nil {
		return fiber.ErrBadRequest
	}

	response, err := c.merchantUsecase.FetchNearby(context, usecase.UserCoordinate{
		Lat: lat,
		Lon: lon,
	}, query)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    response,
	})
}
