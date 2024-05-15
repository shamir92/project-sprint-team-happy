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
}

func NewUserITUsecase(bcryptHelper helper.IBcryptPasswordHash, userRepository repository.IUserRepository, jwtManager helper.IJWTManager) *userITUsecase {
	return &userITUsecase{
		bcryptHelper:   bcryptHelper,
		userRepository: userRepository,
		jwtManager:     jwtManager,
	}
}

type UserITRegisterRequest struct {
	NIP      int    `json:"nip" validate:"required,numeric,min=6150000000000,max=6159999999999"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type UserITRegisterResponse struct {
	Token string `json:"access_token"`
	ID    string `json:"user_id"`
	Name  string `json:"name"`
	NIP   string `json:"nip"`
}

func (u *userITUsecase) RegisterUserIT(userITRequest UserITRegisterRequest) (UserITRegisterResponse, error) {
	var user entity.User

	userNip := strconv.FormatInt(int64(userITRequest.NIP), 10)
	if !entity.ValidateUserNIP(userNip, entity.IT) {
		return UserITRegisterResponse{}, helper.CustomError{
			Message: "NIP is not valid",
			Code:    400,
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
	return UserITRegisterResponse{
		Token: token,
		ID:    user.ID.String(),
		Name:  user.Name,
		NIP:   user.NIP,
	}, nil
}

type UserITLoginRequest struct {
	NIP      int    `json:"nip" validate:"required,numeric,min=6150000000000,max=6159999999999"`
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type UserITLoginResponse struct {
	Token string `json:"access_token"`
	ID    string `json:"user_id"`
	Name  string `json:"name"`
	NIP   string `json:"nip"`
}

func (u *userITUsecase) LoginUserIT(request UserITLoginRequest) (UserITLoginResponse, error) {
	var user entity.User

	userNip := strconv.FormatInt(int64(request.NIP), 10)
	if !entity.ValidateUserNIP(userNip, entity.IT) {
		return UserITLoginResponse{}, helper.CustomError{
			Message: "NIP is not valid",
			Code:    400,
		}
	}
	user, err := u.userRepository.GetByNIP(userNip)
	if err != nil {
		return UserITLoginResponse{}, err
	}
	if !u.bcryptHelper.Compare(user.Password, request.Password) {
		return UserITLoginResponse{}, helper.CustomError{
			Message: "password is not valid",
			Code:    400,
		}
	}

	token, err := u.jwtManager.CreateToken(user)
	if err != nil {
		return UserITLoginResponse{}, err
	}
	return UserITLoginResponse{
		Token: token,
		ID:    user.ID.String(),
		Name:  user.Name,
		NIP:   user.NIP,
	}, nil
}
