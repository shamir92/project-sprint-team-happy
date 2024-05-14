package entity

import (
	"time"

	"github.com/google/uuid"
)

type MedicalRecord struct {
	ID          uuid.UUID `json:"id"`
	PatientID   string    `json:"patient_id"`
	Symptoms    string    `json:"symptoms"`
	Medications string    `json:"medications"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}
