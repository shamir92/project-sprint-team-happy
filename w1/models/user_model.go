package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `db:"id" json:"id"`                 // uuid from github.com/google/uuid
	FullName  string     `db:"fullname" json:"fullname"`     // varchar(50)
	Email     string     `db:"email" json:"email"`           // varchar(50), unique
	Password  string     `db:"password" json:"-"`            // varchar(72), not returned in JSON
	CreatedAt time.Time  `db:"created_at" json:"created_at"` // timestamp with time zone
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"` // timestamp with time zone, nullable
}
