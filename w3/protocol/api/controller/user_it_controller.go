package controller

import (
	"net/http"

	"halosuster/domain/usecase"
	"halosuster/internal/helper"
	"halosuster/protocol/api/dto"

	"github.com/gofiber/fiber/v2"
)

type userITController struct {
	// pingUsecase usecase.IPingUsecase
	userITUsecase usecase.IUserITUsecase
}

// TODO: add all function under ping controller to inferface. this will make it easier to test
type IUserITController interface {
	RegisterUserIT(c *fiber.Ctx) error
	LoginUserIT(c *fiber.Ctx) error
}

func NewUserITController(userITUsecase usecase.IUserITUsecase) *userITController {
	return &userITController{
		userITUsecase: userITUsecase,
	}
}

func (pc *userITController) RegisterUserIT(c *fiber.Ctx) error {
	//TODO: initization variable
	var request usecase.UserITRegisterRequest

	// parse json
	if err := c.BodyParser(&request); err != nil {
		return fiber.ErrBadRequest
	}

	// TODO: Validate Struct
	if err := helper.ValidateStruct(&request); err != nil {
		return err
	}

	data, err := pc.userITUsecase.RegisterUserIT(request)
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(dto.UserITRegisterControllerResponse{
		Message: "user registered successfully",
		Data:    data,
	})
}

func (pc *userITController) LoginUserIT(c *fiber.Ctx) error {
	//TODO: initization variable
	var request usecase.UserITLoginRequest

	// parse json
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// TODO: Validate Struct
	if err := helper.ValidateStruct(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// TODO: add logic
	data, err := pc.userITUsecase.LoginUserIT(request)
	if err != nil {
		// tar ganti
		// dirty dulu.
		return c.Status(http.StatusInternalServerError).JSON(err)
	}
	// TODO: return response

	return c.Status(http.StatusCreated).JSON(dto.UserITRegisterControllerResponse{
		Message: "user login successfully",
		Data:    data,
	})
}
