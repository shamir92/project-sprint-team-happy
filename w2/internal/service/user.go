package service

import (
	"eniqlostore/internal/entity"
	"fmt"
)

type IUserRepository interface {
	Insert(user entity.User) (entity.User, error)
	CheckExistByPhoneNumber(phoneNumber string) (bool, error)
}

type IAuthTokenManager interface {
	CreateToken(user entity.User) (string, error)
}

type IPasswordHash interface {
	Hash(password string) (string, error)
}

type UserServiceDeps struct {
	UserRepository   IUserRepository
	AuthTokenManager IAuthTokenManager
	PasswordHash     IPasswordHash
}

type UserService struct {
	userRepository IUserRepository
	tokenManager   IAuthTokenManager
	passwordHash   IPasswordHash
}

func NewUserService(deps UserServiceDeps) *UserService {
	return &UserService{
		userRepository: deps.UserRepository,
		tokenManager:   deps.AuthTokenManager,
		passwordHash:   deps.PasswordHash,
	}
}

type CreateStaffRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
	Password    string `json:"password"`
}

type CreateUserOut struct {
	entity.User
	AccessToken string `json:"accessToken"`
}

func (s *UserService) UserCreate(in CreateStaffRequest) (CreateUserOut, error) {
	isExist, err := s.userRepository.CheckExistByPhoneNumber(in.PhoneNumber)

	if err != nil {
		return CreateUserOut{}, err
	}

	if isExist {
		return CreateUserOut{}, entity.UserError{
			Message: fmt.Sprintf("user with phone number %s already exist", in.PhoneNumber),
		}
	}

	newUser, err := entity.NewUser(in.PhoneNumber, in.Name, in.Password)

	if err != nil {
		return CreateUserOut{}, err
	}

	hashedPassword, err := s.passwordHash.Hash(newUser.Password)

	if err != nil {
		return CreateUserOut{}, err
	}

	newUser.Password = hashedPassword
	user, err := s.userRepository.Insert(newUser)

	if err != nil {
		return CreateUserOut{}, err
	}

	token, err := s.tokenManager.CreateToken(user)

	if err != nil {
		return CreateUserOut{}, err
	}

	return CreateUserOut{
		User:        user,
		AccessToken: token,
	}, nil
}
