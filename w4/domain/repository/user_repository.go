package repository

import (
	"belimang/domain/entity"
	"database/sql"
	"errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type IUserRepository interface {
	// Return created user's ID
	Insert(newUser entity.User) (entity.User, error)

	// Return true when username is exist in database
	CheckUsernameExist(username string) (bool, error)

	// Find 1 user by its username
	FindOneByUsername(username string) (entity.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Insert(newUser entity.User) (entity.User, error) {
	q := `
		INSERT INTO users(username, password, email, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	createdUser := newUser
	err := r.db.QueryRow(q, newUser.Username, newUser.Password, newUser.Email, newUser.Role).Scan(&createdUser.ID)

	if err != nil {
		return entity.User{}, nil
	}

	return createdUser, nil
}

func (r *userRepository) CheckUsernameExist(username string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT count(username) FROM users WHERE username = $1`, username).Scan(&count)

	return count > 0, err
}

func (r *userRepository) FindOneByUsername(username string) (entity.User, error) {
	q := `SELECT id, username, role, password, email FROM users WHERE username = $1`
	var user entity.User
	err := r.db.QueryRow(q, username).Scan(&user.ID, &user.Username, &user.Role, &user.Password, &user.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserNotFound
		} else {
			return user, err
		}
	}

	return user, nil
}
