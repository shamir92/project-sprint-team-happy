package usecase

import (
	"halosuster/domain/entity"
	"halosuster/internal/helper"
	"time"
)

type iMedicalRecordPatientRepository interface {
	Browse(builder ...entity.BrowseMedicalRecordPatientOptionBuilder) ([]entity.MedicalRecordPatient, error)
	Create(patient entity.MedicalRecordPatient) (entity.MedicalRecordPatient, error)
	Exists(id int) (bool, error)
}

type medicalRecordPatientUsecase struct {
	patientRepository iMedicalRecordPatientRepository
}

type IMedicalRecordPatientUsecase interface {
	Create(req MedicalRecordPatientCreateRequest) (MedicalRecordPatientCreateResponse, error)
	Browse(query MedicalRecordPatientBrowseQuery) ([]MedicalRecordPatientBrowseResponse, error)
}

func NewMedicalRecordPatientUsecase(patientRepository iMedicalRecordPatientRepository) *medicalRecordPatientUsecase {
	return &medicalRecordPatientUsecase{patientRepository}
}

type MedicalRecordPatientCreateRequest struct {
	IdentityNumber      int    `json:"identityNumber" validate:"required,numeric"`
	PhoneNumber         string `json:"phoneNumber" validate:"required,startswith=+62,min=10,max=15"`
	Name                string `json:"name" validate:"required,min=3,max=30"`
	BirthDate           string `json:"birthDate" validate:"required,iso8601"`
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

type MedicalRecordPatientBrowseQuery struct {
	Limit          int    `query:"limit"`
	Offset         int    `query:"offset"`
	IdentityNumber int    `query:"identityNumber"`
	Name           string `query:"name"`
	PhoneNumber    *int   `query:"phoneNumber"`
	SortCreatedAt  string `query:"createdAt"`
}

type MedicalRecordPatientBrowseResponse struct {
	ID          int    `json:"identityNumber"`
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
	BirthDate   string `json:"birthDate"`
	Gender      string `json:"gender"`
	CreatedAt   string `json:"createdAt"`
}

func (u *medicalRecordPatientUsecase) Browse(query MedicalRecordPatientBrowseQuery) ([]MedicalRecordPatientBrowseResponse, error) {
	var (
		options  []entity.BrowseMedicalRecordPatientOptionBuilder
		response []MedicalRecordPatientBrowseResponse = make([]MedicalRecordPatientBrowseResponse, 0)
	)

	if query.Limit <= 0 {
		query.Limit = 5
	}

	if query.Offset < 0 {
		query.Offset = 0
	}

	options = append(options, entity.WithOffsetAndLimit(query.Offset, query.Limit))

	if query.IdentityNumber > 0 {
		options = append(options, entity.WithIdentityNumber(query.IdentityNumber))
	}

	if query.Name != "" {
		options = append(options, entity.WithName(query.Name))
	}

	if query.PhoneNumber != nil {
		options = append(options, entity.WithPhoneNumber(*query.PhoneNumber))
	}

	if query.SortCreatedAt == entity.ASC.String() {
		options = append(options, entity.WithSortCreatedAt(entity.ASC))
	} else {
		options = append(options, entity.WithSortCreatedAt(entity.DESC))
	}

	patients, err := u.patientRepository.Browse(options...)
	if err != nil {
		return response, err
	}

	for _, patient := range patients {
		response = append(response, MedicalRecordPatientBrowseResponse{
			ID:          patient.ID,
			PhoneNumber: patient.PhoneNumber,
			Name:        patient.Name,
			BirthDate:   patient.BirthDate.Format("2006-01-02"),
			Gender:      patient.Gender,
			CreatedAt:   patient.CreatedAt.Format(time.RFC3339),
		})
	}

	return response, nil
}
