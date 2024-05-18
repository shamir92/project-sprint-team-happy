package controller

import (
	"halosuster/domain/usecase"
	"halosuster/internal/helper"
	"halosuster/protocol/api/dto"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type iMedicalRecordPatientUsecase interface {
	Create(req usecase.MedicalRecordPatientCreateRequest) (usecase.MedicalRecordPatientCreateResponse, error)
	Browse(query usecase.MedicalRecordPatientBrowseQuery) ([]usecase.MedicalRecordPatientBrowseResponse, error)
}

type IMedicalRecordPatientController interface {
	Create(c *fiber.Ctx) error
	Browse(c *fiber.Ctx) error
}

type medicalRecordPatientController struct {
	patientUsecase iMedicalRecordPatientUsecase
}

func NewMedicalRecordPatientController(patientUsecase iMedicalRecordPatientUsecase) *medicalRecordPatientController {
	return &medicalRecordPatientController{patientUsecase}
}

func (ctr *medicalRecordPatientController) Create(c *fiber.Ctx) error {
	var request usecase.MedicalRecordPatientCreateRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if err := helper.ValidateStruct(&request); err != nil {
		return err
	}

	data, err := ctr.patientUsecase.Create(request)
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(dto.PatientRegisterControllerResponse{
		Message: "patient registered successfully",
		Data:    data,
	})
}

func (ctr *medicalRecordPatientController) Browse(c *fiber.Ctx) error {
	var query usecase.MedicalRecordPatientBrowseQuery
	c.QueryParser(&query)

	data, err := ctr.patientUsecase.Browse(query)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(dto.PatientRegisterControllerResponse{
		Message: "success",
		Data:    data,
	})
}
