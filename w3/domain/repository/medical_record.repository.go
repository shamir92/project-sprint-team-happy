package repository

import (
	"database/sql"
	"halosuster/domain/entity"
)

type IMedicalRecordRepository interface {
	InsertOne(entity.MedicalRecord) error
}

type medicalRecordRepository struct {
	db *sql.DB
}

func NewMedicalRecordRepository(db *sql.DB) *medicalRecordRepository {
	return &medicalRecordRepository{
		db: db,
	}
}

func (mrr *medicalRecordRepository) InsertOne(newMR entity.MedicalRecord) error {
	q := `
		INSERT INTO 
			medical_records(patient_id, symtomps, medications, created_by) 
		VALUES ($1, $2, $3, $4)
		RETURNING id
		`

	var createdMrID string
	err := mrr.db.QueryRow(q, newMR.PatientID, newMR.Symptoms, newMR.Medications, newMR.CreatedBy).Scan(&createdMrID)

	if err != nil {
		return err
	}

	return nil
}
