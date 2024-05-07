package service

import (
	"eniqlostore/internal/entity"
)

type Service struct {
	userRepository userRepository
}

type userRepository interface {
	Insert(user entity.User) (entity.User, error)
	CheckExistByPhoneNumber(phoneNumber string) (bool, error)
}

type ServiceDeps struct {
	UserRepository userRepository
}

func NewService(opts ServiceDeps) *Service {
	return &Service{
		userRepository: opts.UserRepository,
	}
}
