package controller

import (
	"belimang/domain/usecase"
	"belimang/internal/helper"
	"belimang/protocol/api/dto"
	"context"
	"log"
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
	log.Println(context)

	_, span := tracer.Start(context, "handling register request")
	defer span.End()
	var body usecase.UserRegisterPayload

	if err := ctx.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	if err := helper.ValidateStruct(body); err != nil {
		return err
	}

	token, err := c.userUsecase.Register(body)

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
	var body usecase.UserLoginPayload

	if err := ctx.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	if err := helper.ValidateStruct(body); err != nil {
		return err
	}

	resp, err := c.userUsecase.Login(body)

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
