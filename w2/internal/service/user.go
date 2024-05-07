package service

import (
	"eniqlostore/internal/entity"
	"fmt"
)

type CreateStaffRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
	Password    string `json:"password"`
}

func (s *Service) UserCreate(in CreateStaffRequest) (entity.User, error) {
	isExist, err := s.userRepository.CheckExistByPhoneNumber(in.PhoneNumber)

	if err != nil {
		return entity.User{}, err
	}

	if isExist {
		return entity.User{}, entity.UserError{
			Message: fmt.Sprintf("user with phone number %s already exist", in.PhoneNumber),
		}
	}

	newUser, err := entity.NewUser(in.Password, in.Name, in.Password)

	if err != nil {
		return newUser, err
	}

	user, err := s.userRepository.Insert(newUser)

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
