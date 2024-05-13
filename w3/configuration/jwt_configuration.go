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
	getExpirationTIme() time.Time
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

func (c *jwtConfiguration) getExpirationTIme() time.Time {
	expiresInMinute, _ := strconv.ParseInt(c.expiresInMinute, 10, 64)

	return time.Now().Add(time.Minute * time.Duration(expiresInMinute))
}
