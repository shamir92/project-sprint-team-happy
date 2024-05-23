package configuration

import (
	"os"
	"strconv"
)

type jwtConfiguration struct {
	signingKey      string
	issuer          string
	expiresInMinute string
}

type IJWTConfiguration interface {
	GetSigningKey() string
	GetIssuer() string
	GetExpireInMinute() int
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
	if c.issuer == "" {
		return "app"
	}
	return c.issuer
}

func (c *jwtConfiguration) GetExpireInMinute() int {
	expiresInMinute, err := strconv.ParseInt(c.expiresInMinute, 10, 64)
	if err != nil {
		return 600 // Return the error to be handled by the caller
	}

	return int(expiresInMinute)
}
