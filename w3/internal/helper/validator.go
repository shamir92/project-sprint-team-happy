package helper

import (
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

// Initialize the validator
var validate *validator.Validate

// Initialize the validator in an init function
func init() {
	validate = validator.New()
	validate.RegisterValidation("date", validateDate)
	validate.RegisterValidation("phone", validatePhoneNumber)
	validate.RegisterValidation("iso8601", validateIsISO8601)
	validate.RegisterValidation("nurse_nip", validateNurseNIP)
	validate.RegisterValidation("it_nip", validateITNIP)
}

// Validation function
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func validateDate(fl validator.FieldLevel) bool {
	date, err := time.Parse("2006-01-02", fl.Field().String())

	if date.IsZero() || err != nil {
		return false
	}

	return true
}

func validatePhoneNumber(fl validator.FieldLevel) bool {
	if err := PhoneNumber(fl.Field().String()).Valid(); err != nil {
		return false
	}

	return true
}

func validateIsISO8601(fl validator.FieldLevel) bool {
	date := fl.Field().String()
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05.000Z07:00",
		"2006-01-02",
	}

	for _, format := range formats {
		if _, err := time.Parse(format, date); err == nil {
			return true
		}
	}
	return false
}

func validateNurseNIP(fl validator.FieldLevel) bool {

	// NIP should be an integer represented as a string
	var nurseNip string = strconv.FormatInt(int64(fl.Field().Int()), 10)

	// Rule 1: first three digits should be "303"
	if nurseNip[:3] != "303" {
		return false
	}
	log.Println("masuk sini")

	// Rule 2: fourth digit should be '1' for male, '2' for female
	if nurseNip[3] != '1' && nurseNip[3] != '2' {
		return false
	}

	// Rule 3: fifth to eighth digits should be a year from 2000 to the current year
	year, err := strconv.Atoi(nurseNip[4:8])
	if err != nil {
		return false
	}
	currentYear := time.Now().Year()
	if year < 2000 || year > currentYear {
		return false
	}

	// Rule 4: ninth and tenth digits should be a valid month from "01" to "12"
	month, err := strconv.Atoi(nurseNip[8:10])
	if err != nil {
		return false
	}
	if month < 1 || month > 12 {
		return false
	}

	// Rule 5: eleventh to thirteenth digits should be a valid number from "000" to "99999"
	randomDigits := nurseNip[10:]
	log.Println(randomDigits)
	match, _ := regexp.MatchString(`^\d{3,5}$`, randomDigits)
	if !match {
		return false
	}

	return true
}

func validateITNIP(fl validator.FieldLevel) bool {
	field := fl.Field()

	// NIP should be an integer represented as a string
	nip := strconv.FormatInt(int64(field.Int()), 10)

	// NIP should be exactly 13 digits long
	// if len(nip) != 13 {
	// 	return false
	// }

	// Rule 1: first three digits should be "303"
	if nip[:3] != "615" {
		return false
	}

	// Rule 2: fourth digit should be '1' for male, '2' for female
	if nip[3] != '1' && nip[3] != '2' {
		return false
	}

	// Rule 3: fifth to eighth digits should be a year from 2000 to the current year
	year, err := strconv.Atoi(nip[4:8])
	if err != nil {
		return false
	}
	currentYear := time.Now().Year()
	if year < 2000 || year > currentYear {
		return false
	}

	// Rule 4: ninth and tenth digits should be a valid month from "01" to "12"
	month, err := strconv.Atoi(nip[8:10])
	if err != nil {
		return false
	}
	if month < 1 || month > 12 {
		return false
	}

	// Rule 5: eleventh to thirteenth digits should be a valid number from "000" to "99999"
	randomDigits := nip[10:]
	match, _ := regexp.MatchString(`^\d{3,5}$`, randomDigits)
	if !match {
		return false
	}

	return true
}
