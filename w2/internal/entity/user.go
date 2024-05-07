package entity

import "fmt"

type UserError struct {
	Message string
	Code    int
}

func (u UserError) Error() string {
	return u.Message
}

func (u UserError) HttpStatusCode() int {
	return u.Code
}

type User struct {
	UserID      string `json:"userId"`
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
	Password    string `json:"-"`
}

func validateName(name string) error {
	const MIN_NAME = 5
	const MAX_NAME = 50

	if len(name) < MIN_NAME || len(name) > MAX_NAME {
		return UserError{
			Message: fmt.Sprintf("name: min = %d and max = %d characters", MIN_NAME, MAX_NAME),
		}
	}

	return nil
}

func validatePassword(password string) error {
	const MIN_PASSWORD = 5
	const MAX_PASSWORD = 15

	if len(password) < MIN_PASSWORD || len(password) > MAX_PASSWORD {
		return UserError{
			Message: fmt.Sprintf("password: min = %d and max = %d characters", MIN_PASSWORD, MAX_PASSWORD),
		}
	}

	return nil
}

func NewUser(phoneNumber string, name string, password string) (User, error) {
	var emptyUser User

	// TODO: validate phone number
	if err := validatePassword(password); err != nil {
		return emptyUser, err
	}

	if err := validateName(name); err != nil {
		return emptyUser, err
	}

	return User{
		Name:        name,
		PhoneNumber: phoneNumber,
		Password:    password,
	}, nil
}
