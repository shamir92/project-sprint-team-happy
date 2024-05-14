package configuration

import (
	"os"
	"strconv"
	"time"
)

type jwtConfiguration struct {
	signingKey      string
	issuer          string
	expiresInMinute string
}

type IJWTConfiguration interface {
	GetSigningKey() string
	GetIssuer() string
	GetExpirationTIme() (time.Time, error)
}

func NewJWTConfiguration() *jwtConfiguration {
	return &jwtConfiguration{
		signingKey:      os.Getenv("JWT_SECRET"),
		issuer:          os.Getenv("JWT_ISSUER"),
		expiresInMinute: os.Getenv("JWT_EXPIRES_IN_MINUTE"),
	}
}

func (c *jwtConfiguration) GetSigningKey() string {
	return c.signingKey
}

func (c *jwtConfiguration) GetIssuer() string {
	return c.issuer
}

func (c *jwtConfiguration) GetExpirationTIme() (time.Time, error) {
	expiresInMinute, err := strconv.ParseInt(c.expiresInMinute, 10, 64)
	if err != nil {
		return time.Time{}, err // Return the error to be handled by the caller
	}
	return time.Now().Add(time.Minute * time.Duration(expiresInMinute)), nil
}
