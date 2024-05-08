package auth

import (
	"log"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type bcryptPasswordHash struct {
	saltCost int
}

func NewBcryptPasswordHash() *bcryptPasswordHash {
	salt := os.Getenv("BCRYPT_SALT")

	saltCost, err := strconv.Atoi(salt)

	if err != nil {
		log.Fatalf("bcrypt: salt is not a number")
	}

	return &bcryptPasswordHash{
		saltCost: saltCost,
	}
}

func (b *bcryptPasswordHash) Hash(password string) (hashedPassword string, err error) {
	c := []byte(password)

	bytes, err := bcrypt.GenerateFromPassword(c, b.saltCost)

	hashedPassword = string(bytes)

	return hashedPassword, err
}

func (b *bcryptPasswordHash) Compare(hashedPassword string, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plain))

	return err == nil
}
