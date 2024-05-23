package entity

import "github.com/google/uuid"

type UserRole string

const (
	ROLE_ADMIN UserRole = "ADMIN"
	ROLE_USER  UserRole = "DEFAULT"
)

func (u UserRole) String() string {
	return string(u)
}

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Role     UserRole  `json:"role"`
}
