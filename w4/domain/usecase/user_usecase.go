package usecase

import (
	"belimang/domain/entity"
	"belimang/domain/repository"
	"belimang/internal/helper"
)

type IUserUsecase interface {
	Register(payload UserRegisterPayload) (string, error)
}

type userUsecase struct {
	userRepository repository.IUserRepository
	jwtManager     helper.IJWTManager
	bcrypt         helper.IBcryptPasswordHash
}

func NewUserUsecase(userRepository repository.IUserRepository, jwtManager helper.IJWTManager, bcrypt helper.IBcryptPasswordHash) *userUsecase {
	return &userUsecase{
		userRepository: userRepository,
		bcrypt:         bcrypt,
		jwtManager:     jwtManager,
	}
}

type UserRegisterPayload struct {
	Username string `json:"username" validate:"required,min=5,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=30"`
}

func (u *userUsecase) Register(payload UserRegisterPayload) (string, error) {
	var badTokenValue = ""

	isExist, err := u.userRepository.CheckUsernameExist(payload.Username)
	if err != nil {
		return badTokenValue, err
	}

	if isExist {
		return badTokenValue, helper.CustomError{
			Code:    400,
			Message: "username is already registered",
		}
	}

	hashedPassword, err := u.bcrypt.Hash(payload.Password)
	if err != nil {
		return badTokenValue, err
	}

	insertedUser, err := u.userRepository.Insert(entity.User{
		Email:    payload.Email,
		Username: payload.Username,
		Role:     entity.ROLE_USER,
		Password: hashedPassword,
	})

	if err != nil {
		return badTokenValue, err
	}

	token, err := u.jwtManager.CreateToken(insertedUser)

	if err != nil {
		return badTokenValue, err
	}

	return token, nil
}
