package helper

import (
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
