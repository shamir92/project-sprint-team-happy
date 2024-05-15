package configuration

import (
	"log"
	"os"
	"strconv"
)

type bcryptConfiguration struct {
	bcryptSalt string
}

func NewBcryptConfiguration() *bcryptConfiguration {
	return &bcryptConfiguration{
		bcryptSalt: os.Getenv("BCRYPT_SALT"),
	}
}

type IBcryptConfiguration interface {
	GetBcryptSalt() int
}

func (c *bcryptConfiguration) GetBcryptSalt() int {
	saltCost, err := strconv.Atoi(c.bcryptSalt)

	if err != nil {
		log.Fatalf("bcrypt: salt is not a number")
	}

	return saltCost
}
