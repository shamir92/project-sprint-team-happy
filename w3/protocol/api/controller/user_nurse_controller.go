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
	UpdateNurse(c *fiber.Ctx) error
	DeleteNurse(c *fiber.Ctx) error
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
		return fiber.ErrBadRequest
	}

	if err := helper.ValidateStruct(&request); err != nil {
		return err
	}

	data, err := pc.nurseUsecase.Create(request, user.ID)
	if err != nil {

		return err
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

func (pc *userNurseController) UpdateNurse(c *fiber.Ctx) error {
	var request usecase.UpdateNurseRequest
	if err := c.BodyParser(&request); err != nil {
		return fiber.ErrBadRequest
	}

	if err := helper.ValidateStruct(request); err != nil {
		return err
	}

	nurseUserId := c.Params("userNurseId")
	if err := pc.nurseUsecase.Update(request, nurseUserId); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "nurse updated successfully",
	})
}

func (pc *userNurseController) DeleteNurse(c *fiber.Ctx) error {
	nurseUserId := c.Params("userNurseId")
	if err := pc.nurseUsecase.Delete(nurseUserId); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "nurse deleted successfully",
	})
}

func (pc *userNurseController) SetAccessNurse(c *fiber.Ctx) error {
	var req usecase.SetAccessNurseRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadGateway
	}

	if err := helper.ValidateStruct(req); err != nil {
		return err
	}

	req.UserID = c.Params("userNurseId")
	if err := pc.nurseUsecase.SetAccess(req); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success set access to nurse user",
	})
}

func (pc *userNurseController) LoginNurse(c *fiber.Ctx) error {
	var request usecase.LoginNurseRequest
	if err := c.BodyParser(&request); err != nil {
		return fiber.ErrBadRequest
	}

	if err := helper.ValidateStruct(request); err != nil {
		return err
	}

	data, err := pc.nurseUsecase.Login(request)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "nurse loged in successfully",
		"data": dto.NurseLoginDtoResponse{
			UserID:      data.User.ID.String(),
			NIP:         data.User.NIP,
			Name:        data.User.Name,
			AccessToken: data.AccessToken,
		},
	})
}
