package controller

import (
	"halosuster/domain/usecase"
	"halosuster/internal/helper"

	"github.com/gofiber/fiber/v2"
)

type IMedicalRecordController interface {
	Create(c *fiber.Ctx) error
}

type medicalRecordController struct {
	medicalRecordUsecase usecase.IMedicalRecordUsecase
}

func NewMedicalRecordController(mru usecase.IMedicalRecordUsecase) *medicalRecordController {
	return &medicalRecordController{
		medicalRecordUsecase: mru,
	}
}

func (mrc *medicalRecordController) Create(c *fiber.Ctx) error {
	var body usecase.AddMedicalRecordPayload

	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	if err := helper.ValidateStruct(body); err != nil {
		return err
	}

	user := c.Locals("user").(*helper.JsonWebTokenClaims)

	if err := mrc.medicalRecordUsecase.Create(body, user.UserID); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "success",
	})
}
