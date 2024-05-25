package controller

import (
	"belimang/domain/usecase"
	"belimang/internal/helper"
	"net/http"

	"github.com/gofiber/fiber/v2"
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
	var body usecase.AdminRegisterPayload

	if err := ctx.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	if err := helper.ValidateStruct(body); err != nil {
		return err
	}

	token, err := c.adminUsecase.Register(body)

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
	var body usecase.AdminLoginPayload

	if err := ctx.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	if err := helper.ValidateStruct(body); err != nil {
		return err
	}

	resp, err := c.adminUsecase.Login(body)

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
