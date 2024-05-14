package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID      `json:"id"`
	NIP                 string         `json:"nip"`
	Name                string         `json:"name"`
	Password            sql.NullString `json:"password"`
	Role                string         `json:"role"`
	IdentityCardScanImg string         `json:"identity_card_scan_img"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           sql.NullTime   `json:"deleted_at"`
}
