package models

import (
	"database/sql"
	"errors"
	"fmt"
	"gin-mvc/internal"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	ErrCatNotFound = errors.New("cat not found")
)

type Cat struct {
	ID          uuid.UUID    `db:"id" json:"id"`                     // UUID primary key
	Name        string       `db:"name" json:"name"`                 // VARCHAR(30)
	Sex         string       `db:"sex" json:"sex"`                   // VARCHAR(10)
	AgeInMonth  int          `db:"age_in_month" json:"age_in_month"` // INT
	Description string       `db:"description" json:"description"`   // VARCHAR(20)
	ImageURLs   []string     `db:"image_urls" json:"image_urls"`     // TEXT[], array of strings in Go
	Race        string       `db:"race" json:"race"`                 // VARCHAR(50)
	CreatedBy   uuid.UUID    `db:"created_by" json:"created_by"`     // UUID
	DeletedAt   *time.Time   `db:"deleted_at" json:"deleted_at"`     // TIMESTAMP WITH TIME ZONE, nullable
	HasMatched  bool         `db:"has_matched" json:"has_matched"`   // BOOLEAN with default false
	OwnerID     uuid.UUID    `db:"owner_id" json:"owner_id"`         // UUID
	CreatedAt   time.Time    `db:"created_at" json:"created_at"`     // timestamp with time zone
	UpdatedAt   sql.NullTime `db:"updated_at" json:"updated_at"`     // timestamp with time zone, nullable
}

type CatMatchResponse struct {
	ID          uuid.UUID      `db:"id" json:"id"`                     // UUID primary key
	Name        string         `db:"name" json:"name"`                 // VARCHAR(30)
	Sex         string         `db:"sex" json:"sex"`                   // VARCHAR(10)
	AgeInMonth  int            `db:"age_in_month" json:"age_in_month"` // INT
	Description string         `db:"description" json:"description"`   // VARCHAR(20)
	ImageURLs   pq.StringArray `db:"image_urls" json:"image_urls"`     // TEXT[], array of strings in Go
	Race        string         `db:"race" json:"race"`                 // VARCHAR(50)
	HasMatched  bool           `db:"has_matched" json:"has_matched"`   // BOOLEAN with default false
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`     // timestamp with time zone
}

func (c *Cat) update(in EditCatIn) {
	c.Name = in.Name
	c.Race = in.Race
	c.AgeInMonth = in.Age
	c.Description = in.Description
	c.ImageURLs = in.ImageURLs
	c.UpdatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
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

type EditCatIn struct {
	ID          string
	Name        string
	Race        string
	Age         int
	Description string
	ImageURLs   []string
}

func getCatByIdAndOwnerId(catId string, ownerId string, db *sqlx.DB) (Cat, error) {

	cat := Cat{}

	selectQuery := `
		SELECT 
			id, name, sex, age_in_month, description,
			image_urls, race, owner_id
		FROM 
			cats WHERE id = $1 AND owner_id = $2
	`

	err := db.Get(&cat, selectQuery, catId, ownerId)

	return cat, err
}

func EditCat(in EditCatIn, userId string) (Cat, error) {
	if _, err := uuid.Parse(in.ID); err != nil {
		return Cat{}, CatError{Message: ErrCatNotFound.Error(), Code: http.StatusNotFound}
	}

	db := internal.GetDB()

	cat, err := getCatByIdAndOwnerId(in.ID, userId, db)

	if errors.Is(err, sql.ErrNoRows) {
		return Cat{}, CatError{Message: ErrCatNotFound.Error(), Code: http.StatusNotFound}
	}

	cat.update(in)

	updateQuery := `
		UPDATE 
			cats
		SET 
			name = $1, race = $2, age_in_month = $3, 
			description = $4, image_urls = $5
		WHERE
			id = $6
	`

	res, err := db.Exec(updateQuery, cat.Name, cat.Race, cat.AgeInMonth, cat.Description, pq.Array(cat.ImageURLs), cat.ID)

	if err != nil {
		return Cat{}, err
	}

	if _, err := res.RowsAffected(); err != nil {
		return Cat{}, err
	}

	return cat, nil
}

func CatGetById(id string, db *sqlx.DB) (Cat, error) {
	if db == nil {
		db = internal.GetDB()
	}

	cat := Cat{}
	query := `
		SELECT 
			id, name, sex, owner_id, has_matched
		FROM 
			cats WHERE id = $1 AND deleted_at IS NULL
	`
	err := db.QueryRow(query, id).Scan(&cat.ID, &cat.Name, &cat.Sex, &cat.OwnerID, &cat.HasMatched)

	if errors.Is(err, sql.ErrNoRows) {
		return Cat{}, CatError{Message: ErrCatNotFound.Error(), Code: http.StatusNotFound}
	}

	return cat, err
}
