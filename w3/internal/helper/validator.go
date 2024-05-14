package helper

import "github.com/go-playground/validator/v10"

// Initialize the validator
var validate *validator.Validate

// Initialize the validator in an init function
func init() {
	validate = validator.New()
}

// Validation function
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}
