package auth

import (
	"eniqlostore/internal/entity"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type AuthJwtTokenManager interface {
	CreateToken(user entity.User) (string, error)
	GetClaim(tokenString string) (*JsonWebTokenClaims, error)
}

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
		expirationTimeInMinute: 600,
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
			Issuer:    t.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(t.signingKey))

	return ss, err
}

func (j *jsonWebToken) GetClaim(tokenString string) (*JsonWebTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JsonWebTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		} else if method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}

		return []byte(j.signingKey), nil

	}, jwt.WithExpirationRequired(), jwt.WithIssuer(j.issuer))

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claim, ok := token.Claims.(*JsonWebTokenClaims)

	if !ok {
		return nil, ErrInvalidToken
	}

	return claim, nil
}
