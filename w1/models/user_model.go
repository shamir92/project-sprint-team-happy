package models

import (
	"bytes"
	"database/sql"
	"errors"
	"gin-mvc/internal"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound               = errors.New("user not found")
	ErrUserEmailAlreadyRegistered = errors.New("email already registered")
)

type UserError struct {
	Message    string
	StatusCode int
}

func (u UserError) Error() string {
	return u.Message
}

func (u UserError) HTTPStatusCode() int {
	return u.StatusCode
}

type User struct {
	ID        uuid.UUID  `db:"id" json:"id"`                 // uuid from github.com/google/uuid
	FullName  string     `db:"fullname" json:"fullname"`     // varchar(50)
	Email     string     `db:"email" json:"email"`           // varchar(50), unique
	Password  string     `db:"password" json:"-"`            // varchar(72), not returned in JSON
	CreatedAt time.Time  `db:"created_at" json:"created_at"` // timestamp with time zone
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"` // timestamp with time zone, nullable
}

type RegisterUser struct {
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(newUser RegisterUser) (User, error) {
	db := internal.GetDB()

	// check email already registered or not
	var scannedEmail string
	err := db.QueryRow(`SELECT email FROM users WHERE email = $1`, strings.ToLower(newUser.Email)).Scan(&scannedEmail)

	// if the error is sql.ErrNowRows, it means the email  hasn't registered yet
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Fatalln(err)
		return User{}, err
	}

	if len(scannedEmail) > 0 {
		return User{}, UserError{Message: ErrUserEmailAlreadyRegistered.Error()}
	}

	createdUser := User{
		FullName: newUser.FullName,
		Password: newUser.Password,
		Email:    strings.ToLower(newUser.Email),
	}

	// Bcrypt Salt
	var salt int

	salt, err = strconv.Atoi(os.Getenv("BCRYPT_SALT"))

	if err != nil {
		return User{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createdUser.Password), salt)
	if err != nil {
		return User{}, err
	}

	createdUser.Password = bytes.NewBuffer(hashedPassword).String()

	insertQuery := `
		INSERT INTO users(fullname, email, password)
		VALUES($1, $2, $3) RETURNING id
	`

	err = db.QueryRow(insertQuery, createdUser.FullName, createdUser.Email, createdUser.Password).Scan(&createdUser.ID)

	if err != nil {
		return User{}, err
	}

	return createdUser, nil
}
