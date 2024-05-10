package repository

import (
	"database/sql"
	"eniqlostore/internal/entity"
	"errors"
	"net/http"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db}
}

type IUserRepository interface {
	Insert(user entity.User) (entity.User, error)
	CheckExistByPhoneNumber(phoneNumber string) (bool, error)
	GetByPhoneNumber(phoneNumber string) (entity.User, error)
	GetById(id string) (entity.User, error)
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
		SELECT COUNT(phone_number) FROM users WHERE phone_number = $1
	`

	var count int
	err := r.db.QueryRow(query, phoneNumber).Scan(&count)

	// A users by given phone number  exist when the scanned phone number isn't empty string
	return count > 0, err
}

func (r *userRepository) GetByPhoneNumber(phoneNumber string) (entity.User, error) {
	query := `
		SELECT user_id, name, phone_number, password FROM users WHERE phone_number = $1
	`

	var user entity.User
	err := r.db.QueryRow(query, phoneNumber).Scan(&user.UserID, &user.Name, &user.PhoneNumber, &user.Password)

	// User is not registered in db
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return user, entity.UserError{
			Message: "user not found",
			Code:    http.StatusNotFound,
		}
	}

	// Unknown error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) GetById(id string) (entity.User, error) {
	query := `
		SELECT user_id, name, phone_number, password FROM users WHERE user_id = $1
	`
	var user entity.User
	err := r.db.QueryRow(query, id).Scan(&user.UserID, &user.Name, &user.PhoneNumber, &user.Password)
	if err != nil {
		return user, err
	}

	return user, err
}
