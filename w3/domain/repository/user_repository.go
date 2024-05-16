package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"halosuster/domain/entity"
	"halosuster/internal/helper"
	"log"
	"net/http"
	"time"
)

type userRepository struct {
	db *sql.DB
}

var (
	errUserNotFound = errors.New("user not found")
	errUpdateUser   = errors.New("update user failed")
)

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db}
}

type IUserRepository interface {
	GetByNIP(nip string) (entity.User, error)
	InsertUser(user entity.User) (entity.User, error)
	CheckNIPExist(nip string) (bool, error)
	GetUserNurseByID(userId string) (entity.User, error)
	Update(entity.User) error
	Delete(userId string) error
	UpdatePassword(userId string, newHashedPassword string) error
}

func (r *userRepository) GetByNIP(nip string) (entity.User, error) {
	query := `
		SELECT id,  nip, name, password, role FROM users WHERE nip = $1
	`

	var user entity.User
	err := r.db.QueryRow(query, nip).Scan(&user.ID, &user.NIP, &user.Name, &user.Password, &user.Role)

	// User is not registered in db
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return user, helper.CustomError{
			Message: "data not found",
			Code:    http.StatusNotFound,
		}
	}

	return user, nil
}

func (r *userRepository) InsertUser(user entity.User) (entity.User, error) {
	query := `
		INSERT INTO users(nip, name, password, role) 
		VALUES($1, $2, $3, $4) 
		ON CONFLICT DO NOTHING
		RETURNING id, nip
	`

	err := r.db.QueryRow(query, user.NIP, user.Name, user.Password, user.Role).Scan(&user.ID, &user.NIP)
	if err != nil {
		// err is not nil if the user is already registered
		if !errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, helper.CustomError{
				Message: "Duplicate User",
				Code:    400,
			}
		}
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) CheckNIPExist(nip string) (bool, error) {
	query := `
		SELECT count(nip) password FROM users WHERE nip = $1 AND deleted_at IS NULL
	`

	var count int
	err := r.db.QueryRow(query, nip).Scan(&count)

	return count > 0, err
}

func (r *userRepository) GetUserNurseByID(userId string) (entity.User, error) {
	query := `
		SELECT id,  nip, name, role FROM users WHERE id = $1 AND role = $2
	`

	var nurse entity.User

	err := r.db.QueryRow(query, userId, entity.NURSE).Scan(&nurse.ID, &nurse.NIP, &nurse.Name, &nurse.Role)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, helper.CustomError{
				Code:    404,
				Message: errUserNotFound.Error(),
			}
		} else {
			log.Printf("GetUserNurseByID: %v", err)
			return entity.User{}, err
		}

	}

	return nurse, nil
}

func (r *userRepository) Update(user entity.User) error {
	query := `UPDATE users SET name = $1, nip = $2 WHERE id = $3`

	res, err := r.db.Exec(query, user.Name, user.NIP, user.ID.String())

	if err != nil {
		log.Printf("failed to update user: %v => user: %v", err, user)
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Printf("failed to update user: %v", err)
		return err
	}

	// Rows affected should be one
	if rowsAffected != 1 {
		log.Printf("failed to update user: rows affected greater than 1 - error: %v", err)
		return errUpdateUser
	}

	return nil
}

func (r *userRepository) Delete(userId string) error {
	query := `UPDATE users SET nip = NULL, deleted_at = $1 WHERE id = $2`

	res, err := r.db.Exec(query, time.Now(), userId)

	if err != nil {
		log.Printf("failed to delete user: %v => user: %s", err, userId)
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Printf("failed to delete user: %v", err)
		return err
	}

	// Rows affected should be one
	if rowsAffected != 1 {
		log.Printf("failed to delete user: rows affected greater than 1 - error: %v", err)
		return fmt.Errorf("failed to delete user %s", userId)
	}

	return nil
}

func (r *userRepository) UpdatePassword(userId string, newHashedPassword string) error {
	query := `UPDATE users SET password = $1 WHERE id = $2`

	res, err := r.db.Exec(query, newHashedPassword, userId)

	if err != nil {
		log.Printf("failed to update user's password: %v => user: %s", err, userId)
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Printf("failed to update user's password: %v", err)
		return err
	}

	// Rows affected should be one
	if rowsAffected != 1 {
		log.Printf("failed to update user's password: rows affected greater than 1 - error: %v", err)
		return errUpdateUser
	}

	return nil
}
