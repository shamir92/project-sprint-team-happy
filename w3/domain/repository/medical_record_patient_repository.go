package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"halosuster/domain/entity"
	"strings"
)

type medicalRecordPatientRepository struct {
	db *sql.DB
}

type IMedicalRecordPatientRepository interface {
	Exists(id int) (bool, error)
	Create(patient entity.MedicalRecordPatient) (entity.MedicalRecordPatient, error)
	Browse(builder ...entity.BrowseMedicalRecordPatientOptionBuilder) ([]entity.MedicalRecordPatient, error)
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

func (r *medicalRecordPatientRepository) Browse(builder ...entity.BrowseMedicalRecordPatientOptionBuilder) ([]entity.MedicalRecordPatient, error) {
	var (
		patients   []entity.MedicalRecordPatient
		conditions []string = make([]string, 0)
		values     []any    = make([]any, 0)
		sort       string
	)

	options := &entity.BrowseMedicalRecordPatientOption{}
	for _, o := range builder {
		o(options)
	}

	if options.IdentityNumber != nil {
		values = append(values, *options.IdentityNumber)
		conditions = append(conditions, fmt.Sprintf("id = $%d", len(values)))
	}

	if options.Name != "" {
		values = append(values, "%"+options.Name+"%")
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", len(values)))
	}

	if options.PhoneNumber != nil {
		values = append(values, fmt.Sprintf("+%d%%", *options.PhoneNumber))
		conditions = append(conditions, fmt.Sprintf("phone_number LIKE $%d", len(values)))
	}

	if options.Limit == 0 {
		options.Limit = 5
	}

	if options.SortCreatedAt.String() != "" {
		sort = fmt.Sprintf("ORDER BY created_at %s", options.SortCreatedAt.String())
	}

	WHERE := ""
	if len(conditions) > 0 {
		WHERE = "WHERE"
	}

	query := fmt.Sprintf(`
		SELECT id, phone_number, name, birth_date, gender, created_at FROM public.patients
		%s %s %s LIMIT %d OFFSET %d
	`, WHERE, strings.Join(conditions, " AND "), sort, options.Limit, options.Offset)

	fmt.Println(query, values)
	res, err := r.db.Query(query, values...)
	if err != nil {
		return patients, err
	}
	for res.Next() {
		var patient entity.MedicalRecordPatient
		err = res.Scan(&patient.ID, &patient.PhoneNumber, &patient.Name, &patient.BirthDate, &patient.Gender, &patient.CreatedAt)
		if err != nil {
			return []entity.MedicalRecordPatient{}, err
		}

		patients = append(patients, patient)
	}

	return patients, nil
}
