package entity

import "time"

type Patient struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	PhoneNumber         string    `json:"phone_number"`
	BirthDate           time.Time `json:"birth_date"`
	Gender              string    `json:"gender"`
	IdentityCardScanImg string    `json:"identity_card_scan_img"`
	CreatedAt           time.Time `json:"created_at"`
}
