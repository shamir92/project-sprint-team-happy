package usecase

import (
	"halosuster/domain/entity"
	"halosuster/domain/repository"
	"halosuster/internal/helper"
	"strconv"
)

type userITUsecase struct {
	bcryptHelper   helper.IBcryptPasswordHash
	userRepository repository.IUserRepository
	jwtManager     helper.IJWTManager
}

type IUserITUsecase interface {
	RegisterUserIT(userITRequest UserITRegisterRequest) (UserITRegisterResponse, error)
	LoginUserIT(request UserITLoginRequest) (UserITLoginResponse, error)
	GetUsers(in entity.ListUserPayload) ([]entity.User, error)
}

func NewUserITUsecase(bcryptHelper helper.IBcryptPasswordHash, userRepository repository.IUserRepository, jwtManager helper.IJWTManager) *userITUsecase {
	return &userITUsecase{
		bcryptHelper:   bcryptHelper,
		userRepository: userRepository,
		jwtManager:     jwtManager,
	}
}

type UserITRegisterRequest struct {
	NIP      int    `json:"nip" validate:"required,numeric,it_nip"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type UserITRegisterResponse struct {
	Token string `json:"accessToken"`
	ID    string `json:"userId"`
	Name  string `json:"name"`
	NIP   int    `json:"nip"`
}

func (u *userITUsecase) RegisterUserIT(userITRequest UserITRegisterRequest) (UserITRegisterResponse, error) {
	var user entity.User

	userNip := strconv.FormatInt(int64(userITRequest.NIP), 10)

	userNipExists, err := u.userRepository.CheckNIPExist(userNip)
	if err != nil {
		return UserITRegisterResponse{}, err
	}

	if userNipExists {
		return UserITRegisterResponse{}, helper.CustomError{
			Message: "NIP already exists",
			Code:    409,
		}
	}

	hashedPassword, err := u.bcryptHelper.Hash(userITRequest.Password)
	if err != nil {
		return UserITRegisterResponse{}, helper.CustomError{
			Message: err.Error(),
			Code:    500,
		}
	}
	user.Name = userITRequest.Name

	user.NIP = userNip
	user.Password = hashedPassword
	user.Role = string(entity.IT)

	user, err = u.userRepository.InsertUser(user)

	if err != nil {
		return UserITRegisterResponse{}, err
	}

	token, err := u.jwtManager.CreateToken(user)
	if err != nil {
		return UserITRegisterResponse{}, err
	}
	integer, _ := strconv.Atoi(user.NIP)
	return UserITRegisterResponse{
		Token: token,
		ID:    user.ID.String(),
		Name:  user.Name,
		NIP:   integer,
	}, nil
}

type UserITLoginRequest struct {
	NIP      int    `json:"nip" validate:"required,numeric,it_nip"`
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type UserITLoginResponse struct {
	Token string `json:"accessToken"`
	ID    string `json:"userId"`
	Name  string `json:"name"`
	NIP   int    `json:"nip"`
}

func (u *userITUsecase) LoginUserIT(request UserITLoginRequest) (UserITLoginResponse, error) {
	var user entity.User

	userNip := strconv.FormatInt(int64(request.NIP), 10)

	user, err := u.userRepository.GetByNIP(userNip)
	if err != nil {
		return UserITLoginResponse{}, helper.CustomError{
			Message: "nip or password doesn't match",
			Code:    404,
		}
	}
	if !u.bcryptHelper.Compare(user.Password, request.Password) {
		return UserITLoginResponse{}, helper.CustomError{
			Message: "nip or password doesn't match",
			Code:    400,
		}
	}

	token, err := u.jwtManager.CreateToken(user)
	if err != nil {
		return UserITLoginResponse{}, err
	}

	integer, _ := strconv.Atoi(user.NIP)
	return UserITLoginResponse{
		Token: token,
		ID:    user.ID.String(),
		Name:  user.Name,
		NIP:   integer,
	}, nil
}

func (u *userITUsecase) GetUsers(in entity.ListUserPayload) ([]entity.User, error) {
	return u.userRepository.List(in)
}
