package controller

import (
	"net/http"

	"halosuster/domain/usecase"
	"halosuster/internal/helper"
	"halosuster/protocol/api/dto"

	"github.com/gofiber/fiber/v2"
)

type userNurseController struct {
	nurseUsecase usecase.IUserNurseUsecase
}

// TODO: add all function under ping controller to inferface. this will make it easier to test
type IUserNurseController interface {
	CreateNurse(c *fiber.Ctx) error
}

func NewUserNurseController(nurseUsecase usecase.IUserNurseUsecase) *userNurseController {
	return &userNurseController{
		nurseUsecase: nurseUsecase,
	}
}

func (pc *userNurseController) CreateNurse(c *fiber.Ctx) error {
	user := c.Locals("user").(*helper.JsonWebTokenClaims)

	var request usecase.CreateNurseRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if err := helper.ValidateStruct(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	data, err := pc.nurseUsecase.Create(request, user.ID)
	if err != nil {
		// tar ganti
		// dirty dulu.
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "user registered successfully",
		"data": dto.CreateUserNurseDtoResponse{
			Name:   data.Name,
			NIP:    data.NIP,
			UserID: data.ID.String(),
		},
	})
}
