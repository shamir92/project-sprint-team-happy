package entity

import (
	"database/sql"

	"github.com/google/uuid"
)

type MedicalRecord struct {
	ID          int          `json:"id"`
	PatientID   string       `json:"patient_id"`
	Symptoms    string       `json:"symptoms"`
	Medications string       `json:"medications"`
	CreatedBy   uuid.UUID    `json:"created_by"`
	CreatedAt   sql.NullTime `json:"created_at"`
}
