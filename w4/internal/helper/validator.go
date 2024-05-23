package helper

import (
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

// Initialize the validator
var validate *validator.Validate

// Initialize the validator in an init function
func init() {
	validate = validator.New()
	validate.RegisterValidation("date", validateDate)
	validate.RegisterValidation("iso8601", validateIsISO8601)
	validate.RegisterValidation("urlformat", validateURLFormat)

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

const (
	URLSchema    = `((ftp|tcp|udp|wss?|https?):\/\/)`
	URLUsername  = `(\S+(:\S*)?@)`
	URLPath      = `((\/|\?|#)[^\s]*)`
	URLPort      = `(:(\d{1,5}))`
	URLIP        = `([1-9]\d?|1\d\d|2[01]\d|22[0-3]|24\d|25[0-5])(\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-5])){3}`
	URLSubdomain = `((www\.)|([a-zA-Z0-9]+([-_\.]?[a-zA-Z0-9])*[a-zA-Z0-9]\.[a-zA-Z0-9]+))`
	URLRegex     = `^` + URLSchema + `?` + URLUsername + `?` + `((` + URLIP + `|\[IPv6:(.*)\]|` + URLSubdomain + `))` + URLPort + `?` + URLPath + `?$`
)

func validateURLFormat(fl validator.FieldLevel) bool {
	identityCardScanImg := fl.Field().String()
	regex := regexp.MustCompile(URLRegex)
	return regex.MatchString(identityCardScanImg)
}
