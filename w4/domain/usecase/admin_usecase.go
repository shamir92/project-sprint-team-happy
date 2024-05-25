package usecase

import (
	"belimang/domain/entity"
	"belimang/domain/repository"
	"belimang/internal/helper"
	"errors"
	"log"
)

var (
	ErrLoginAdminFailed = errors.New("username or password is wrong")
)

type IAdminUsecase interface {
	Register(payload AdminRegisterPayload) (string, error)
	Login(payload AdminLoginPayload) (AdminLoginResp, error)
}

type adminUsecase struct {
	userRepository repository.IUserRepository
	jwtManager     helper.IJWTManager
	bcrypt         helper.IBcryptPasswordHash
}

func NewAdminUsecase(userRepository repository.IUserRepository, jwtManager helper.IJWTManager, bcrypt helper.IBcryptPasswordHash) *adminUsecase {
	return &adminUsecase{
		userRepository: userRepository,
		bcrypt:         bcrypt,
		jwtManager:     jwtManager,
	}
}

type AdminRegisterPayload struct {
	Username string `json:"username" validate:"required,min=5,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=30"`
}

type AdminLoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     entity.UserRole
}

type AdminLoginResp struct {
	Token string
}

func (u *adminUsecase) Register(payload AdminRegisterPayload) (string, error) {
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
		Role:     entity.ROLE_ADMIN,
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

func (u *adminUsecase) Login(payload AdminLoginPayload) (AdminLoginResp, error) {
	user, err := u.userRepository.FindOneByUsername(payload.Username)

	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return AdminLoginResp{}, helper.CustomError{
				Message: ErrLoginAdminFailed.Error(),
				Code:    400,
			}
		}

		log.Fatalf("LOGIN_ADMIN | FindOneByUsername | error: %v\n", err)
		return AdminLoginResp{}, err
	}

	if !user.IsAdmin() {
		return AdminLoginResp{}, helper.CustomError{
			Message: ErrLoginAdminFailed.Error(),
			Code:    400,
		}
	}

	if isMatch := u.bcrypt.Compare(user.Password, payload.Password); !isMatch {
		return AdminLoginResp{}, helper.CustomError{
			Message: ErrLoginAdminFailed.Error(),
			Code:    400,
		}
	}

	token, err := u.jwtManager.CreateToken(user)

	if err != nil {
		log.Fatalf("LOGIN_ADMIN | fail when creating the token | error: %v\n", err)
		return AdminLoginResp{}, err
	}

	return AdminLoginResp{
		Token: token,
	}, err
}
