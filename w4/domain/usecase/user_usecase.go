package usecase

import (
	"belimang/domain/entity"
	"belimang/domain/repository"
	"belimang/internal/helper"
	"errors"
	"log"
)

var (
	ErrLoginUserFailed = errors.New("username or password is wrong")
)

type IUserUsecase interface {
	Register(payload UserRegisterPayload) (string, error)
	Login(payload UserLoginPayload) (UserLoginResp, error)
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

type UserLoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResp struct {
	Token string
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

func (u *userUsecase) Login(payload UserLoginPayload) (UserLoginResp, error) {
	user, err := u.userRepository.FindOneByUsername(payload.Username)

	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return UserLoginResp{}, helper.CustomError{
				Message: ErrLoginUserFailed.Error(),
				Code:    400,
			}
		}

		log.Printf("LOGIN | FindOneByUsername | error: %v\n", err)
		return UserLoginResp{}, err
	}

	if !user.IsUserRole() {
		log.Printf("LOGIN | IsUserRole | %s", user.Role)
		return UserLoginResp{}, helper.CustomError{
			Message: ErrLoginUserFailed.Error(),
			Code:    400,
		}
	}

	if isMatch := u.bcrypt.Compare(user.Password, payload.Password); !isMatch {
		return UserLoginResp{}, helper.CustomError{
			Message: ErrLoginUserFailed.Error(),
			Code:    400,
		}
	}

	token, err := u.jwtManager.CreateToken(user)

	if err != nil {
		log.Fatalf("LOGIN | fail when creating the token | error: %v\n", err)
		return UserLoginResp{}, err
	}

	return UserLoginResp{
		Token: token,
	}, err
}
