package controller

import (
	"halosuster/domain/entity"
	"halosuster/domain/usecase"
	"halosuster/internal/helper"
	"halosuster/protocol/api/dto"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type IMedicalRecordController interface {
	Create(c *fiber.Ctx) error
	GetMedicalRecords(c *fiber.Ctx) error
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

func (mrc *medicalRecordController) GetMedicalRecords(c *fiber.Ctx) error {
	var query entity.ListMedicalRecordPayload

	err := c.QueryParser(&query)

	if err != nil {
		log.Printf("GetMedicalRecords: failed to parse query - err: %v\n", err)
	}

	medicalRecords, err := mrc.medicalRecordUsecase.GetRecords(query)

	if err != nil {
		return err
	}

	var medicalRecordsResponseData []dto.MedicalRecordItemDto

	for _, mc := range medicalRecords {
		patient := mc.GetPatient()
		user := mc.GetUserCreatedBy()
		nip, _ := strconv.Atoi(user.NIP)

		medicalRecordsResponseData = append(medicalRecordsResponseData, dto.MedicalRecordItemDto{
			IdentityDetail: dto.IdentityDetailDto{
				IdentityNumber:      int64(patient.ID),
				PhoneNumber:         patient.PhoneNumber,
				Name:                patient.Name,
				BirthDate:           patient.BirthDate.Format(time.RFC3339),
				Gender:              patient.Gender,
				IdentityCardScanImg: patient.IdentityCardScanImg,
			},
			Symptoms:    mc.Symptoms,
			Medications: mc.Medications,
			CreatedAt:   mc.CreatedAt.Format(time.RFC3339),
			CreatedBy: dto.CreatedByDto{
				Name:   user.Name,
				Nip:    int64(nip),
				UserId: user.ID.String(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    medicalRecordsResponseData,
	})
}
