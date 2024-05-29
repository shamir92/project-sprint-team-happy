package controller

import (
	"belimang/domain/usecase"
	"belimang/internal/helper"
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

type adminController struct {
	adminUsecase usecase.IAdminUsecase
}

func NewAdminController(adminUsecase usecase.IAdminUsecase) *adminController {
	return &adminController{
		adminUsecase: adminUsecase,
	}
}

func (c *adminController) Register(ctx *fiber.Ctx) error {
	context := ctx.Locals("ctx").(context.Context)
	tracer := ctx.Locals("tracer").(trace.Tracer)

	_, span := tracer.Start(context, "Register")
	defer span.End()
	var body usecase.AdminRegisterPayload

	if err := ctx.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	if err := helper.ValidateStruct(body); err != nil {
		return err
	}

	token, err := c.adminUsecase.Register(context, body)

	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"data": fiber.Map{
			"token": token,
		},
	})
}

func (c *adminController) Login(ctx *fiber.Ctx) error {
	context := ctx.Locals("ctx").(context.Context)
	tracer := ctx.Locals("tracer").(trace.Tracer)

	_, span := tracer.Start(context, "Login")
	defer span.End()
	var body usecase.AdminLoginPayload

	if err := ctx.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	if err := helper.ValidateStruct(body); err != nil {
		return err
	}

	resp, err := c.adminUsecase.Login(context, body)

	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"data": fiber.Map{
			"token": resp.Token,
		},
	})
}
