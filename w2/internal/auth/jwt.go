package auth

import (
	"eniqlostore/internal/entity"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jsonWebToken struct {
	expirationTimeInMinute int
	signingKey             string
	issuer                 string
}

type JsonWebTokenClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func NewJwt() *jsonWebToken {
	return &jsonWebToken{
		expirationTimeInMinute: 60,
		signingKey:             os.Getenv("JWT_SECRET"),
		issuer:                 "app",
	}
}

func (t *jsonWebToken) CreateToken(user entity.User) (string, error) {
	expiresAt := time.Now().Add(time.Duration(t.expirationTimeInMinute) * time.Minute)

	claims := JsonWebTokenClaims{
		user.UserID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(t.signingKey))

	return ss, err
}
