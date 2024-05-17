package usecase

import (
	"errors"
	"halosuster/domain/entity"
	"halosuster/domain/repository"
	"halosuster/internal/helper"
	"log"
	"strconv"

	"github.com/google/uuid"
)

var (
	ErrIdentityNumberNotFound = errors.New("identityNumber is not found")
)

type IMedicalRecordUsecase interface {
	Create(in AddMedicalRecordPayload, createdBy string) error
	GetRecords(in entity.ListMedicalRecordPayload) ([]entity.MedicalRecord, error)
}

type medicalRecordUsecase struct {
	patientRepository       repository.IMedicalRecordPatientRepository
	medicalRecordRepository repository.IMedicalRecordRepository
}

func NewMedicalRecordUsecase(mrpr repository.IMedicalRecordPatientRepository, mrr repository.IMedicalRecordRepository) *medicalRecordUsecase {
	return &medicalRecordUsecase{
		patientRepository:       mrpr,
		medicalRecordRepository: mrr,
	}
}

type AddMedicalRecordPayload struct {
	IdentityNumber int    `json:"identityNumber" validate:"required,numeric"`
	Symptoms       string `json:"symptoms" validate:"required,min=1,max=2000"`
	Medications    string `json:"medications" validate:"required,min=1,max=2000"`
}

func (mru *medicalRecordUsecase) Create(in AddMedicalRecordPayload, createdBy string) error {
	isExist, err := mru.patientRepository.Exists(in.IdentityNumber)

	if err != nil {
		return err
	}

	if !isExist {
		return helper.CustomError{
			Code:    400,
			Message: ErrIdentityNumberNotFound.Error(),
		}
	}

	createdByUser, err := uuid.Parse(createdBy)

	if err != nil {
		log.Printf("failed to create new medical record: %v\n", err)
		return err
	}

	newMR := entity.MedicalRecord{
		PatientID:   strconv.FormatInt(int64(in.IdentityNumber), 10),
		Symptoms:    in.Symptoms,
		Medications: in.Medications,
		CreatedBy:   createdByUser,
	}

	return mru.medicalRecordRepository.InsertOne(newMR)
}

func (mru *medicalRecordUsecase) GetRecords(in entity.ListMedicalRecordPayload) ([]entity.MedicalRecord, error) {
	return mru.medicalRecordRepository.Find(in)
}
