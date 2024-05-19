package entity

import (
	"database/sql"
	"log"
	"net/url"
	"path/filepath"
	"regexp"
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
	IdentityCardScanImg string       `json:"identityCardScanImg"`
	CreatedAt           time.Time    `json:"createdAt"`
	UpdatedAt           time.Time    `json:"updatedAt"`
	DeletedAt           sql.NullTime `json:"deletedAt"`
}

func (u User) IsNurse() bool {
	return u.Role == string(NURSE)
}

func (u User) HasAccess() bool {
	return u.Password != ""
}

// true when the NIP is valid
func ValidateUserNIP(nip string, role UserRole) bool {
	// Check length
	// as per the latest requirement, nip length changed to 15 (from 13)
	if len(nip) != 15 {
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
	// as per the latest requirement, nip length changed to 15 (from 13)
	// if _, err := strconv.Atoi(nip[10:15]); err != nil {
	// 	fmt.Println("disini?", err.Error())

	// 	return false
	// }

	randomDigits := nip[10:]
	match, _ := regexp.MatchString(`^\d{3,5}$`, randomDigits)
	if !match {
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

func IsValidUserRole(rawRole string) bool {
	role := UserRole(strings.ToUpper(rawRole))
	return NURSE == role || IT == role
}

type ListUserPayload struct {
	UserID          string `query:"userId"`
	Limit           string `query:"limit"`
	Offset          string `query:"offset"`
	Name            string `query:"name"`
	NIP             int    `query:"nip"`
	Role            string `query:"role"`
	SortByCreatedAt string `query:"createdAt"`
}
