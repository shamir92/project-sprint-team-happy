package repository

import (
	"belimang/domain/entity"
	"database/sql"
)

type IUserRepository interface {
	// Return created user's ID
	Insert(newUser entity.User) (entity.User, error)

	// Return true when username is exist in database
	CheckUsernameExist(username string) (bool, error)
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
