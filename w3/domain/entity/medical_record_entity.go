package entity

import (
	"time"

	"github.com/google/uuid"
)

type MedicalRecord struct {
	ID            int       `json:"id"`
	PatientID     string    `json:"patient_id"`
	Symptoms      string    `json:"symptoms"`
	Medications   string    `json:"medications"`
	CreatedBy     uuid.UUID `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	patient       MedicalRecordPatient
	userCreatedBy User
}

func (m *MedicalRecord) SetPatient(patient MedicalRecordPatient) {
	m.patient = patient
}

func (m *MedicalRecord) SetUserCreatedBy(user User) {
	m.userCreatedBy = user
}

func (m *MedicalRecord) GetPatient() MedicalRecordPatient {
	return m.patient
}

func (m *MedicalRecord) GetUserCreatedBy() User {
	return m.userCreatedBy
}

type ListMedicalRecordPayload struct {
	IdentityNumber  int    `query:"identityDetail.identityNumber"`
	CreatedByUserId string `query:"createdBy.userId"`
	CreatedByNip    string `query:"createdBy.nip"`
	Limit           string `query:"limit"`
	Offset          string `query:"offset"`
	SortByCreatedAt string `query:"createdAt"`
}
