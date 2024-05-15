package helper

import (
	"halosuster/configuration"

	"golang.org/x/crypto/bcrypt"
)

type bcryptPasswordHash struct {
	saltCost int
}

type IBcryptPasswordHash interface {
	Hash(password string) (hashedPassword string, err error)
	Compare(hashedPassword string, plain string) bool
	GetSaltCost() int
}

func NewBcryptPasswordHash(bcrypt configuration.IBcryptConfiguration) *bcryptPasswordHash {
	bcryptSalt := bcrypt.GetBcryptSalt()
	return &bcryptPasswordHash{
		saltCost: bcryptSalt,
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

func (b *bcryptPasswordHash) GetSaltCost() int {
	return b.saltCost
}
