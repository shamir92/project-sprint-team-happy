package repository

import (
	"database/sql"
	"errors"
	"halosuster/domain/entity"
)

type medicalRecordPatientRepository struct {
	db *sql.DB
}

type IMedicalRecordPatientRepository interface {
	Exists(id int) (bool, error)
	Create(patient entity.MedicalRecordPatient) (entity.MedicalRecordPatient, error)
}

func NewMedicalRecordPatientRepository(db *sql.DB) *medicalRecordPatientRepository {
	return &medicalRecordPatientRepository{db}
}

func (r *medicalRecordPatientRepository) Create(patient entity.MedicalRecordPatient) (entity.MedicalRecordPatient, error) {
	statement := `
		INSERT INTO public.patients(id, name, phone_number, birth_date, gender, identity_card_scan_img) VALUES ($1, $2, $3, $4, $5, $6) RETURNING created_at
	`

	err := r.db.QueryRow(statement, patient.ID, patient.Name, patient.PhoneNumber, patient.BirthDate, patient.Gender, patient.IdentityCardScanImg).Scan(&patient.CreatedAt)
	if err != nil {
		return entity.MedicalRecordPatient{}, err
	}

	return patient, nil
}

func (r *medicalRecordPatientRepository) Exists(id int) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM public.patients WHERE id = $1)`

	err := r.db.QueryRow(query, id).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return exists, err
	}

	return exists, nil
}
