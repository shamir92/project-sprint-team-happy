package repository

import (
	"database/sql"
	"fmt"
	"halosuster/domain/entity"
	"strconv"

	"github.com/google/uuid"
)

type IMedicalRecordRepository interface {
	InsertOne(entity.MedicalRecord) error
	Find(entity.ListMedicalRecordPayload) ([]entity.MedicalRecord, error)
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

func (mrr *medicalRecordRepository) Find(query entity.ListMedicalRecordPayload) ([]entity.MedicalRecord, error) {
	q := `
		SELECT
			mc.id, mc.symtomps, mc.medications, mc.created_at,
			p.id as patient_id, p.name as patient_name, p.birth_date patient_birth_date,
			p.gender as patient_gender, p.identity_card_scan_img AS patient_identity,
			u.name as user_created_name, u.nip as user_created_nip, u.id as user_created_id
		FROM medical_records mc
		INNER JOIN patients p ON p.id = mc.patient_id
		INNER JOIN users u ON u.id = mc.created_by
	`

	var paramsCounter = 1
	var params []interface{}

	if patientId := strconv.FormatInt(int64(query.IdentityNumber), 10); patientId != "" && query.IdentityNumber != 0 {
		q += whereOrAnd(paramsCounter)
		q += fmt.Sprintf("mc.patient_id = $%d", paramsCounter)
		paramsCounter += 1
		params = append(params, patientId)
	}

	if _, err := uuid.Parse(query.CreatedByUserId); err == nil {
		q += fmt.Sprintf("%s u.id = $%d", whereOrAnd(paramsCounter), paramsCounter)
		paramsCounter += 1
		params = append(params, query.CreatedByUserId)
	}

	if query.CreatedByNip != "" {
		q += fmt.Sprintf("%s u.nip = $%d", whereOrAnd(paramsCounter), paramsCounter)
		paramsCounter += 1
		params = append(params, query.CreatedByNip)
	}

	if query.SortByCreatedAt == "asc" || query.SortByCreatedAt == "desc" {
		q += fmt.Sprintf(" ORDER BY created_at %s", query.SortByCreatedAt)
	} else {
		q += " ORDER BY created_at DESC "
	}

	// LIMIT AND OFFSET
	var limit, offset = 5, 0

	if l, err := strconv.Atoi(query.Limit); err == nil && l > 0 {
		limit = l
	}

	if o, err := strconv.Atoi(query.Offset); err == nil && o > 0 {
		offset = o
	}

	q += fmt.Sprintf(" LIMIT %d OFFSET %d ", limit, offset)

	rows, err := mrr.db.Query(q, params...)

	if err != nil {
		return []entity.MedicalRecord{}, err
	}

	defer rows.Close()

	var records []entity.MedicalRecord
	for rows.Next() {
		var mr entity.MedicalRecord
		var user entity.User
		var patient entity.MedicalRecordPatient

		err := rows.Scan(&mr.ID,
			&mr.Symptoms,
			&mr.Medications,
			&mr.CreatedAt,
			&patient.ID,
			&patient.Name,
			&patient.BirthDate,
			&patient.Gender,
			&patient.IdentityCardScanImg,
			&user.Name,
			&user.NIP,
			&user.ID,
		)

		if err != nil {
			return []entity.MedicalRecord{}, err
		}

		mr.SetPatient(patient)
		mr.SetUserCreatedBy(user)

		records = append(records, mr)
	}

	return records, nil
}
