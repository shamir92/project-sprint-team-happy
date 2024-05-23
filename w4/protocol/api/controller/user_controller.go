package controller

import (
	"belimang/domain/usecase"
	"belimang/internal/helper"
	"belimang/protocol/api/dto"
	"net/http"

	"github.com/gofiber/fiber/v2"
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
