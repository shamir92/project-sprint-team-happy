package entity

import (
	"database/sql"
	"log"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	IT    UserRole = "IT"
	NURSE UserRole = "NURSE"
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

// true when the NIP is valid
func ValidateUserNIP(nip string, role UserRole) bool {
	// Check length
	if len(nip) != 13 {
		return false
	}

	// Check role
	if role == IT && nip[:3] != "615" {
		return false
	}

	if role == NURSE && nip[:3] != "303" {
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

func ValidateIdentityCardScanImageURL(rawUrl string) bool {
	if !strings.HasPrefix(rawUrl, "http://") && !strings.HasPrefix(rawUrl, "https://") {
		return false
	}

	u, err := url.ParseRequestURI(rawUrl)

	if err != nil {
		log.Printf("identity card scan image url: %v", err)
		return false
	}

	ext := strings.ToLower(filepath.Ext(u.Path))

	validExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
	}

	_, ok := validExtensions[ext]

	return ok
}
