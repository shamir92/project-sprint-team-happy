package usecase

import (
	"halosuster/domain/entity"
	"halosuster/internal/helper"
	"time"
)

type iPatientRepository interface {
	Create(patient entity.MedicalRecordPatient) (entity.MedicalRecordPatient, error)
	Exists(id int) (bool, error)
}

type medicalRecordPatientUsecase struct {
	patientRepository iPatientRepository
}

type IMedicalRecordPatientUsecase interface {
	Create(req MedicalRecordPatientCreateRequest) (MedicalRecordPatientCreateResponse, error)
}

func NewMedicalRecordPatientUsecase(patientRepository iPatientRepository) *medicalRecordPatientUsecase {
	return &medicalRecordPatientUsecase{patientRepository}
}

type MedicalRecordPatientCreateRequest struct {
	IdentityNumber      int    `json:"identityNumber" validate:"required,numeric"`
	PhoneNumber         string `json:"phoneNumber" validate:"required,phone,min=10,max=15"`
	Name                string `json:"name" validate:"required,min=3,max=30"`
	BirthDate           string `json:"birthDate" validate:"required,date"`
	Gender              string `json:"gender" validate:"required,oneof=male female"`
	IdentityCardScanImg string `json:"identityCardScanImg" validate:"required"`
}

type MedicalRecordPatientCreateResponse struct {
	IdentityNumber int    `json:"identityNumber"`
	Name           string `json:"name"`
	CreatedAt      string `json:"createdAt"`
}

func (u *medicalRecordPatientUsecase) checkDuplicateIdentityNumber(identityNumber int) error {
	exist, err := u.patientRepository.Exists(identityNumber)
	if err != nil {
		return err
	}

	if exist {
		return helper.CustomError{
			Message: "identityNumber already exists",
			Code:    409,
		}
	}

	return nil
}

func (u *medicalRecordPatientUsecase) Create(req MedicalRecordPatientCreateRequest) (MedicalRecordPatientCreateResponse, error) {
	var patient entity.MedicalRecordPatient

	if err := u.checkDuplicateIdentityNumber(req.IdentityNumber); err != nil {
		return MedicalRecordPatientCreateResponse{}, err
	}

	birthDate, _ := time.Parse("2006-01-02", req.BirthDate)

	patient.ID = req.IdentityNumber
	patient.PhoneNumber = req.PhoneNumber
	patient.Name = req.Name
	patient.BirthDate = birthDate
	patient.Gender = req.Gender
	patient.IdentityCardScanImg = req.IdentityCardScanImg

	patient, err := u.patientRepository.Create(patient)
	if err != nil {
		return MedicalRecordPatientCreateResponse{}, err
	}

	return MedicalRecordPatientCreateResponse{
		IdentityNumber: patient.ID,
		Name:           patient.Name,
		CreatedAt:      patient.CreatedAt.Format(time.RFC3339),
	}, nil
}
