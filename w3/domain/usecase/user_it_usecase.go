package usecase

import (
	"halosuster/domain/entity"
	"halosuster/domain/repository"
	"halosuster/internal/helper"
)

type userITUsecase struct {
	bcryptHelper   helper.IBcryptPasswordHash
	userRepository repository.IUserRepository
}

type IUserITUsecase interface {
	RegisterUserIT(userITRequest UserITRegisterRequest) error
}

func NewUserITUsecase(bcryptHelper helper.IBcryptPasswordHash, userRepository repository.IUserRepository) *userITUsecase {
	return &userITUsecase{
		bcryptHelper:   bcryptHelper,
		userRepository: userRepository,
	}
}

type UserITRegisterRequest struct {
	NIP      string `json:"nip" validate:"required, numeric, min=6150000000000, max=6159999999999"`
	Name     string `json:"name" validate:"required, min=5, max=50"`
	Password string `json:"password" validate:"required, min=5, max=33"`
}

func (u *userITUsecase) RegisterUserIT(userITRequest UserITRegisterRequest) error {
	var user entity.User
	if !user.ValidateNIP(userITRequest.NIP) {
		return helper.CustomError{
			Message: "NIP is not valid",
			Code:    400,
		}
	}
	hashedPassword, err := u.bcryptHelper.Hash(userITRequest.Password)
	if err != nil {
		return helper.CustomError{
			Message: err.Error(),
			Code:    500,
		}
	}
	user.Name = userITRequest.Name
	user.NIP = userITRequest.NIP
	user.Password = hashedPassword
	if err := u.userRepository.InsertUserIT(user); err != nil {
		return err
	}
	return nil
}
