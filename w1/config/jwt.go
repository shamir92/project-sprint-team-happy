package config

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jsonToken struct {
	expireInMinute string
	signingKey     string
	issuer         string
}

func NewJWT() *jsonToken {
	return &jsonToken{
		signingKey:     os.Getenv("JWT_SECRET"),
		issuer:         os.Getenv("JWT_ISSUER"),
		expireInMinute: os.Getenv("JWT_EXPIRES_IN_MINUTE"),
	}
}

// bingung mau di pakai tanpa inject interface. hahaha
type IJsonToken interface {
	GetSigningKey() string
	GetIssuer() string
	GetExpireInMinute() (int, error)
}

func (jt *jsonToken) GetSigningKey() string {
	return jt.signingKey
}

func (jt *jsonToken) GetIssuer() string {
	return jt.issuer
}

func (jt *jsonToken) getExpirationTIme() time.Time {
	expiresInMinute, _ := strconv.Atoi(os.Getenv("JWT_EXPIRES_IN_MINUTE"))
	expirationTime := time.Now().Add(time.Duration(expiresInMinute) * time.Minute)
	return expirationTime
}

func (jt *jsonToken) GenerateJWT(data map[string]interface{}) (string, error) {
	// Set up the JWT claims with the provided data
	claims := jwt.MapClaims{}
	for key, value := range data {
		claims[key] = value
	}
	claims["nbf"] = time.Now().Unix()
	claims["exp"] = jt.getExpirationTIme()

	// Create a new JWT token with the claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(jt.signingKey))
}
