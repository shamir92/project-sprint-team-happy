package models

import (
	"time"

	"github.com/google/uuid"
)

type Cat struct {
	ID          uuid.UUID  `db:"id" json:"id"`                     // UUID primary key
	Name        string     `db:"name" json:"name"`                 // VARCHAR(30)
	Sex         string     `db:"sex" json:"sex"`                   // VARCHAR(10)
	AgeInMonth  int        `db:"age_in_month" json:"age_in_month"` // INT
	Description string     `db:"description" json:"description"`   // VARCHAR(20)
	ImageURLs   []string   `db:"image_urls" json:"image_urls"`     // TEXT[], array of strings in Go
	Race        string     `db:"race" json:"race"`                 // VARCHAR(50)
	CreatedBy   uuid.UUID  `db:"created_by" json:"created_by"`     // UUID
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`     // TIMESTAMP WITH TIME ZONE, nullable
	HasMatched  bool       `db:"has_matched" json:"has_matched"`   // BOOLEAN with default false
	OwnerID     uuid.UUID  `db:"owner_id" json:"owner_id"`         // UUID
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`     // timestamp with time zone
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`     // timestamp with time zone, nullable
}
