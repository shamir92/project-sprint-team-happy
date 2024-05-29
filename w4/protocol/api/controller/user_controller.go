package controller

import (
	"belimang/domain/usecase"
	"belimang/internal/helper"
	"belimang/protocol/api/dto"
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

type userController struct {
	userUsecase usecase.IUserUsecase
}

func NewUserController(userUsecase usecase.IUserUsecase) *userController {
	return &userController{
		userUsecase: userUsecase,
	}
}

func (c *userController) Register(ctx *fiber.Ctx) error {
	context := ctx.Locals("ctx").(context.Context)
	tracer := ctx.Locals("tracer").(trace.Tracer)

	_, span := tracer.Start(context, "Register")
	defer span.End()
	var body usecase.UserRegisterPayload

	if err := ctx.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	if err := helper.ValidateStruct(body); err != nil {
		return err
	}

	token, err := c.userUsecase.Register(context, body)

	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"data": dto.UserRegisterDtoResponse{
			Token: token,
		},
	})
}

func (c *userController) Login(ctx *fiber.Ctx) error {
	context := ctx.Locals("ctx").(context.Context)
	tracer := ctx.Locals("tracer").(trace.Tracer)
	_, span := tracer.Start(context, "Login")
	defer span.End()

	var body usecase.UserLoginPayload

	if err := ctx.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	if err := helper.ValidateStruct(body); err != nil {
		return err
	}

	resp, err := c.userUsecase.Login(context, body)

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
