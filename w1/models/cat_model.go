package models

import (
	"database/sql"
	"errors"
	"fmt"
	"gin-mvc/internal"
	"net/http"
	"strconv"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	ErrCatNotFound   = errors.New("cat not found")
	ErrCatIsNotOwned = errors.New("cat is not owned")
)

type Cat struct {
	ID             uuid.UUID      `db:"id" json:"id"`                     // UUID primary key
	Name           string         `db:"name" json:"name"`                 // VARCHAR(30)
	Sex            string         `db:"sex" json:"sex"`                   // VARCHAR(10)
	AgeInMonth     int            `db:"age_in_month" json:"age_in_month"` // INT
	Description    string         `db:"description" json:"description"`   // VARCHAR(20)
	ImageURLs      pq.StringArray `db:"image_urls" json:"image_urls"`     // TEXT[], array of strings in Go
	Race           string         `db:"race" json:"race"`                 // VARCHAR(50)
	CreatedBy      uuid.UUID      `db:"created_by" json:"created_by"`     // UUID
	DeletedAt      sql.NullTime   `db:"deleted_at" json:"deleted_at"`     // TIMESTAMP WITH TIME ZONE, nullable
	HasMatched     bool           `db:"has_matched" json:"has_matched"`   // BOOLEAN with default false
	OwnerID        uuid.UUID      `db:"owner_id" json:"owner_id"`         // UUID
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`     // timestamp with time zone
	UpdatedAt      sql.NullTime   `db:"updated_at" json:"updated_at"`     // timestamp with time zone, nullable
	HasMatchedWith sql.NullBool
}

type CatError struct {
	Message string
	Code    int
}

func (c CatError) Error() string {
	return fmt.Sprintf(c.Message)
}

func (c CatError) HTTPStatusCode() int {
	return c.Code
}

type CreateOrUpdateCatIn struct {
	ID          string
	Name        string
	Race        string
	Sex         string
	Age         int
	Description string
	ImageURLs   []string
}

type GetCatOption struct {
	ID         string `form:"id"`
	Limit      *int   `form:"limit, omitempty"`
	Offset     *int   `form:"offset, omitempty"`
	Race       string `form:"race"`
	Sex        string `form:"sex"`
	Age        string `form:"ageInMonth"`
	Owned      string `form:"owned"` // true | false
	Search     string `form:"search"`
	HasMatched string `form:"hasMatched"` // true | false
}

/*
*

	Parsing ageInMonth string:

	- `ageInMonth=>4` searches data that have more than 4 months
	- `ageInMonth=<4` searches data that have less than 4 months
	- `ageInMonth=4` searches data that have exact 4 month

*
*/
func (opt GetCatOption) ParseAge() (op string, ageInMonth int, valid bool) {

	var (
		lt  = "<"
		gt  = ">"
		lte = "<="
		gte = ">="
		eq  = "="
	)

	opWithValue := opt.Age

	// The minimum length of op+age is 2: (e.g: =1..9)
	if len(opWithValue) < 2 {
		return op, ageInMonth, valid
	}

	// Find the index where the operator ends and the age part begins
	var ageStartIdx int
	for i, char := range opWithValue {
		if unicode.IsDigit(char) {
			ageStartIdx = i
			break
		}
	}

	// Extract the operator and age part based on the index found
	opPart := opWithValue[:ageStartIdx]
	agePart := opWithValue[ageStartIdx:]

	age, err := strconv.Atoi(string(agePart))

	if err != nil {
		return op, ageInMonth, valid // default value
	}

	switch opPart {
	case lt, gt, eq, lte, gte:
		op = opPart
		ageInMonth = age
		valid = true
	}

	return
}

type CatOut struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Sex         string         `json:"sex"`
	AgeInMonth  int            `json:"ageInMonth"`
	Description string         `json:"description"`
	ImageURLs   pq.StringArray `json:"imageUrls"`
	Race        string         `json:"race"`
	HasMatched  bool           `json:"hasMatched"`
	OwnerID     uuid.UUID      `json:"ownerId"`
	CreatedAt   time.Time      `json:"createdAt"`
}

func GetCayById(catId string, db *sqlx.DB) (Cat, error) {
	cat := Cat{}

	selectQuery := `
		SELECT 
			id, name, sex, age_in_month, description,
			image_urls, race, owner_id, has_matched
		FROM 
			cats WHERE id = $1 AND deleted_at IS NULL
	`

	err := db.Get(&cat, selectQuery, catId)

	return cat, err
}

func GetCatByIdAndOwnerId(catId string, ownerId string, db *sqlx.DB) (Cat, error) {
	cat := Cat{}

	selectQuery := `
		SELECT 
			id, name, sex, age_in_month, description,
			image_urls, race, owner_id
		FROM 
			cats WHERE id = $1 AND owner_id = $2 AND deleted_at IS NULL
	`

	err := db.Get(&cat, selectQuery, catId, ownerId)

	return cat, err
}

func CreateCat(in CreateOrUpdateCatIn, userId string) (Cat, error) {
	if !IsValidCatRace(in.Race) {
		return Cat{}, CatError{
			Message: "race is not valid",
			Code:    http.StatusBadRequest,
		}
	}

	db := internal.GetDB()
	insertCatQuery := `
		INSERT INTO cats (name, race, sex, age_in_month, description, image_urls, owner_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at;
	`

	createdCat := Cat{}
	err := db.QueryRow(insertCatQuery, in.Name, in.Race, in.Sex, in.Age, in.Description, pq.Array(in.ImageURLs), userId).Scan(&createdCat.ID, &createdCat.CreatedAt)
	if err != nil {
		return Cat{}, err
	}

	return createdCat, nil
}

// TODO: impelement edit cat's sex requirement
func EditCat(in CreateOrUpdateCatIn, userId string) error {
	if _, err := uuid.Parse(in.ID); err != nil {
		return CatError{Message: ErrCatNotFound.Error(), Code: http.StatusNotFound}
	}

	if !IsValidCatRace(in.Race) {
		return CatError{
			Message: "race is not valid",
			Code:    http.StatusBadRequest,
		}
	}

	db := internal.GetDB()

	cat, err := GetCayById(in.ID, db)

	if err != nil {
		return CatError{Message: ErrCatNotFound.Error(), Code: http.StatusNotFound}
	}

	if cat.HasMatched {
		return CatError{
			Message: "can't update sex when the cat have already matched",
			Code:    http.StatusBadRequest,
		}
	}

	updateQuery := `
		UPDATE 
			cats
		SET 
			name = $1, race = $2, age_in_month = $3, 
			description = $4, image_urls = $5, updated_at = $6, 
			sex = $7
		WHERE
			(id = $8 AND owner_id = $9)
	`

	res, err := db.Exec(updateQuery, in.Name, in.Race, in.Age, in.Description, pq.Array(in.ImageURLs), time.Now(), in.Sex, in.ID, userId)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected < 1 {
		return CatError{Message: ErrCatNotFound.Error(), Code: http.StatusNotFound}

	}

	return nil
}

// TODO: validate cat based on requirement matches before deleting the cat
func DeleteCatById(catId string, userId string) error {
	db := internal.GetDB()
	if _, err := uuid.Parse(catId); err != nil {
		return CatError{Message: ErrCatNotFound.Error(), Code: http.StatusNotFound}
	}

	cat, err := GetCatByIdAndOwnerId(catId, userId, db)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return CatError{Message: ErrCatNotFound.Error(), Code: http.StatusNotFound}
	}

	if cat.HasMatched {
		return CatError{Message: "cat can't be deleted when already matched", Code: http.StatusBadRequest}
	}

	softDeleteCatQuery := `
		UPDATE cats
		SET deleted_at = $1
		WHERE (id = $2 AND owner_id = $3)
	`

	result, err := db.Exec(softDeleteCatQuery, time.Now(), catId, userId)

	if err != nil {
		return err
	}

	totalAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if totalAffected < 1 {
		return CatError{Message: ErrCatNotFound.Error(), Code: http.StatusNotFound}
	}

	return nil
}

func CheckIfCatIsOwnedByUser(ownerId string, userId string) error {
	if ownerId != userId {
		return CatError{Message: ErrCatIsNotOwned.Error(), Code: http.StatusNotFound}
	}

	return nil
}

func GetCats(opts GetCatOption, userId string) ([]CatOut, error) {
	db := internal.GetDB()

	values := []interface{}{}

	query := `
		SELECT 
			id, name, sex,
			age_in_month, description,
			image_urls, race, owner_id,
			created_at
		FROM cats
		WHERE deleted_at IS NULL
	`

	if _, err := uuid.Parse(opts.ID); err == nil {
		values = append(values, opts.ID)
		query += fmt.Sprintf(" AND id = $%d", len(values))
	}

	if owned, err := strconv.ParseBool(opts.Owned); err == nil {
		if owned {
			values = append(values, userId)
			query += fmt.Sprintf(" AND owner_id = $%d", len(values))
		}

		if !owned {
			values = append(values, userId)
			query += fmt.Sprintf(" AND owner_id != $%d", len(values))
		}

	}

	if IsValidCatRace(opts.Race) {
		values = append(values, opts.Race)
		query += fmt.Sprintf(" AND race = $%d", len(values))
	}

	if IsValidCatSex(opts.Sex) {
		values = append(values, opts.Sex)
		query += fmt.Sprintf(" AND sex = $%d", len(values))
	}

	if _, err := strconv.ParseBool(opts.HasMatched); err == nil && opts.HasMatched != "" {
		values = append(values, opts.HasMatched)
		query += fmt.Sprintf(" AND has_matched = $%d", len(values))
	}

	if opts.Search != "" {
		values = append(values, fmt.Sprintf("%%%s%%", opts.Search))
		query += fmt.Sprintf(" AND name ILIKE $%d", len(values))
	}

	if op, age, valid := opts.ParseAge(); valid {
		query += fmt.Sprintf(" AND age_in_month %s %d", op, age)
	}

	offset, limit := 0, 5

	if opts.Offset != nil {
		offset = *opts.Offset
	}

	if opts.Limit != nil {
		limit = *opts.Limit
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT %d OFFSET %d", limit, offset)

	rows, err := db.Query(query, values...)

	if err != nil {
		return []CatOut{}, err
	}

	defer rows.Close()

	var cats []CatOut = []CatOut{}
	for rows.Next() {
		var c CatOut

		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Sex,
			&c.AgeInMonth,
			&c.Description,
			&c.ImageURLs,
			&c.Race,
			&c.OwnerID,
			&c.CreatedAt,
		); err != nil {
			return []CatOut{}, err
		}

		cats = append(cats, c)
	}

	if err := rows.Err(); err != nil {
		return cats, err
	}

	return cats, nil
}

func UpdateHasMatchedCat(catIds []string, db *sqlx.DB, hasMatched bool) error {
	arg := map[string]interface{}{
		"catIds":     catIds,
		"hasMatched": hasMatched,
	}
	query := `UPDATE cats SET has_matched = :hasMatched WHERE id IN (:catIds)`
	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}

	query = db.Rebind(query)
	_, err = db.Exec(query, args...)

	return err
}
