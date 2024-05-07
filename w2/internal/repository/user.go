package repository

import (
	"database/sql"
	"eniqlostore/internal/entity"
	"errors"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) Insert(user entity.User) (entity.User, error) {
	query := `
		INSERT INTO users(phone_number, name, password) 
		VALUES($1, $2, $3) 
		RETURNING user_id`

	newUser := user

	err := r.db.QueryRow(query, user.PhoneNumber, user.Name, user.Password).Scan(&newUser.UserID)

	if err != nil {
		return entity.User{}, err
	}

	return newUser, nil
}

func (r *userRepository) CheckExistByPhoneNumber(phoneNumber string) (bool, error) {
	query := `
		SELECT phone_number FROM users WHERE phone_number = $1
	`

	var scannedPhoneNumber string = ""
	err := r.db.QueryRow(query, phoneNumber).Scan(&scannedPhoneNumber)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil // No rows means user with given phone number isn't exist in database
		} else {
			return false, err
		}
	}

	// A users by given phone number  exist when the scanned phone number isn't empty string
	return scannedPhoneNumber != "", nil
}
