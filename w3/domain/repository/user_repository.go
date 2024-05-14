package repository

import (
	"database/sql"
	"errors"
	"halosuster/domain/entity"
	"halosuster/internal/helper"
	"net/http"
)

type userRepository struct {
	db *sql.DB
}

var (
	errUserNotFound = errors.New("user not found")
)

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db}
}

type IUserRepository interface {
	GetByNIP(nip string) (entity.User, error)
	InsertUserIT(user entity.User) error
}

func (r *userRepository) GetByNIP(nip string) (entity.User, error) {
	query := `
		SELECT nip, name, password FROM users WHERE nip = $1
	`

	var user entity.User
	err := r.db.QueryRow(query, nip).Scan(&user.NIP, &user.Name, &user.Password)

	// User is not registered in db
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return user, helper.CustomError{
			Message: "data not found",
			Code:    http.StatusNotFound,
		}
	}

	return user, nil
}

func (r *userRepository) InsertUserIT(user entity.User) error {
	query := `
		INSERT INTO users(nip, name, password) 
		VALUES($1, $2, $3) 
		RETURNING nip`

	err := r.db.QueryRow(query, user.NIP, user.Name, user.Password).Scan(&user.NIP)
	if err != nil {
		return err
	}

	return nil
}
