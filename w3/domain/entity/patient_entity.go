package entity

import (
	"time"

	"github.com/google/uuid"
)

type Patient struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"name"`
	PhoneNumber         string    `json:"phone_number"`
	BirthDate           time.Time `json:"birth_date"`
	Gender              string    `json:"gender"`
	IdentityCardScanImg string    `json:"identity_card_scan_img"`
	CreatedAt           time.Time `json:"created_at"`
}
