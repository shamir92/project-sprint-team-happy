package repository

import (
	"belimang/domain/entity"
	"context"
	"database/sql"
	"errors"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type IUserRepository interface {
	// Return created user's ID
	Insert(ctx context.Context, newUser entity.User) (entity.User, error)

	// Return true when username is exist in database
	CheckUsernameExist(ctx context.Context, username string) (bool, error)

	// Find 1 user by its username
	FindOneByUsername(ctx context.Context, username string) (entity.User, error)
}

type userRepository struct {
	db     *sql.DB
	tracer trace.Tracer
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		db:     db,
		tracer: otel.Tracer("user-repository"),
	}
}

func (r *userRepository) Insert(ctx context.Context, newUser entity.User) (entity.User, error) {
	_, span := r.tracer.Start(ctx, "Insert")
	defer span.End()

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

func (r *userRepository) CheckUsernameExist(ctx context.Context, username string) (bool, error) {
	_, span := r.tracer.Start(ctx, "CheckUsernameExist")
	defer span.End()

	var count int
	err := r.db.QueryRow(`SELECT count(username) FROM users WHERE username = $1`, username).Scan(&count)

	return count > 0, err
}

func (r *userRepository) FindOneByUsername(ctx context.Context, username string) (entity.User, error) {
	_, span := r.tracer.Start(ctx, "FindOneByUsername")
	defer span.End()

	q := `SELECT id, username, role, password, email FROM users WHERE username = $1`
	var user entity.User
	err := r.db.QueryRow(q, username).Scan(&user.ID, &user.Username, &user.Role, &user.Password, &user.Email)
	log.Println(user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserNotFound
		} else {
			return user, err
		}
	}

	return user, nil
}
