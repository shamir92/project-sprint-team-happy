package entity

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID    `json:"id"`
	NIP                 string       `json:"nip"`
	Name                string       `json:"name"`
	Password            string       `json:"password"`
	Role                string       `json:"role"`
	IdentityCardScanImg string       `json:"identity_card_scan_img"`
	CreatedAt           time.Time    `json:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at"`
	DeletedAt           sql.NullTime `json:"deleted_at"`
}

func (e *User) ValidateNIP(nip string) bool {
	// Check length
	if len(nip) != 13 {
		return false
	}

	// Check prefix
	if nip[:3] != "615" {
		return false
	}

	// Check fourth digit for gender
	if nip[3] != '1' && nip[3] != '2' {
		return false
	}

	// Check year (5th to 8th digit)
	year, err := strconv.Atoi(nip[4:8])
	if err != nil {
		return false
	}
	currentYear := time.Now().Year()
	if year < 2000 || year > currentYear {
		return false
	}

	// Check month (9th and 10th digit)
	month, err := strconv.Atoi(nip[8:10])
	if err != nil {
		return false
	}
	if month < 1 || month > 12 {
		return false
	}

	// Check random digits (11th to 13th digit)
	if _, err := strconv.Atoi(nip[10:13]); err != nil {
		return false
	}

	return true
}
