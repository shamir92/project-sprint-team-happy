package usecase

import (
	"errors"
	"halosuster/domain/entity"
	"halosuster/domain/repository"
	"halosuster/internal/helper"
	"log"
	"strconv"

	"github.com/google/uuid"
)

var (
	errInvalidNIP    = errors.New("NIP is invalid")
	errConflictNIP   = errors.New("NIP is already registered")
	errNurseNotFound = errors.New("nurse not found")
	errNurseAccess   = errors.New("nurse doesn't have access yet")
	errNurseAuth     = errors.New("password is wrong")
)

type IUserNurseUsecase interface {
	Create(in CreateNurseRequest, createdBy string) (entity.User, error)
	Update(in UpdateNurseRequest, nurseUserId string) error
	Delete(nurseId string) error
	SetAccess(SetAccessNurseRequest) error
	Login(in LoginNurseRequest) (LoginNurseResponse, error)
}

type userNurseUseCase struct {
	userRepository repository.IUserRepository
	bcryptHelper   helper.IBcryptPasswordHash
	jwtManager     helper.IJWTManager
}

func NewUserNurseUseCase(userRepo repository.IUserRepository, bcryptHelper helper.IBcryptPasswordHash, jwtManager helper.IJWTManager) *userNurseUseCase {
	return &userNurseUseCase{
		userRepository: userRepo,
		bcryptHelper:   bcryptHelper,
		jwtManager:     jwtManager,
	}
}

type CreateNurseRequest struct {
	NIP                 int    `json:"nip" validate:"required,numeric,nurse_nip"`
	Name                string `json:"name" validate:"required,min=5,max=50"`
	IndetityCardScanImg string `json:"identityCardScanImg" validate:"required,url"`
}

type UpdateNurseRequest struct {
	NIP  int    `json:"nip" validate:"required,numeric,nurse_nip"`
	Name string `json:"name" validate:"required,min=5,max=50"`
}

type SetAccessNurseRequest struct {
	UserID   string
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type LoginNurseRequest struct {
	NIP      int    `json:"nip" validate:"required,numeric,nurse_nip"`
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type LoginNurseResponse struct {
	entity.User
	AccessToken string
}

func (u *userNurseUseCase) Create(in CreateNurseRequest, createdBy string) (entity.User, error) {
	userNip := strconv.FormatInt(int64(in.NIP), 10)
	// userNip := in.NIP
	newNurse := entity.User{
		NIP:                 userNip,
		Name:                in.Name,
		IdentityCardScanImg: in.IndetityCardScanImg,
		Role:                string(entity.NURSE),
	}

	if !entity.ValidateUserNIP(userNip, entity.NURSE) {
		return entity.User{}, helper.CustomError{
			Message: "nip is invalid",
			Code:    400,
		}
	}

	if !entity.ValidateIdentityCardScanImageURL(newNurse.IdentityCardScanImg) {
		return entity.User{}, helper.CustomError{
			Message: "identityCardScanImgUrl is invalid",
			Code:    400,
		}
	}

	isNIPExist, err := u.userRepository.CheckNIPExist(newNurse.NIP)

	if err != nil {
		return entity.User{}, err
	}

	if isNIPExist {
		return entity.User{}, helper.CustomError{
			Message: "nip is already used",
			Code:    409,
		}
	}

	if createdNurse, err := u.userRepository.InsertUser(newNurse); err != nil {
		return entity.User{}, err
	} else {
		newNurse.ID = createdNurse.ID
	}

	return newNurse, nil
}

func (u *userNurseUseCase) Update(in UpdateNurseRequest, nurseUserId string) error {
	nip := strconv.FormatInt(int64(in.NIP), 10)

	// if !entity.ValidateUserNIP(nip, entity.NURSE) {
	// 	return helper.CustomError{
	// 		Code:    400,
	// 		Message: errInvalidNIP.Error(),
	// 	}
	// }
	// nip := in.NIP

	nurse, err := u.getUserNurseByID(nurseUserId)

	if err != nil {
		return err
	}

	if nip != nurse.NIP {
		isExist, err := u.userRepository.CheckNIPExist(nip)

		if err != nil {
			return err
		}

		if isExist {
			return helper.CustomError{
				Code:    409,
				Message: errConflictNIP.Error(),
			}
		}
	}

	nurse.NIP = nip
	nurse.Name = in.Name

	if err := u.userRepository.Update(nurse); err != nil {
		return err
	}

	return nil
}

func (u *userNurseUseCase) Delete(nurseId string) error {
	_, err := u.getUserNurseByID(nurseId)

	if err != nil {
		return err
	}

	if err := u.userRepository.Delete(nurseId); err != nil {
		return err
	}

	return nil
}

func (u *userNurseUseCase) SetAccess(in SetAccessNurseRequest) error {
	nurse, err := u.getUserNurseByID(in.UserID)

	if err != nil {
		return err
	}

	hashedPassword, err := u.bcryptHelper.Hash(in.Password)

	if err != nil {
		log.Printf("SetAccess: %v", err)
		return err
	}

	nurse.Password = hashedPassword

	err = u.userRepository.UpdatePassword(nurse.ID.String(), nurse.Password)

	if err != nil {
		return err
	}

	return nil
}

func (u *userNurseUseCase) Login(in LoginNurseRequest) (LoginNurseResponse, error) {
	var empty LoginNurseResponse
	nip := strconv.FormatInt(int64(in.NIP), 10)
	// log.Println(nip)
	// nip := in.NIP
	user, err := u.userRepository.GetByNIP(nip)

	if err != nil {
		return empty, err
	}

	if !user.IsNurse() {
		return empty, helper.CustomError{
			Code:    404,
			Message: errNurseNotFound.Error(),
		}
	}

	if !user.HasAccess() {
		return empty, helper.CustomError{
			Code:    400,
			Message: errNurseAccess.Error(),
		}
	}

	if !u.bcryptHelper.Compare(user.Password, in.Password) {
		return empty, helper.CustomError{
			Code:    400,
			Message: errNurseAuth.Error(),
		}
	}

	token, err := u.jwtManager.CreateToken(user)

	if err != nil {
		return empty, err
	}

	return LoginNurseResponse{
		User:        user,
		AccessToken: token,
	}, nil
}

func isValidUUID(userId string) bool {
	err := uuid.Validate(userId)

	return err == nil
}

func (u *userNurseUseCase) getUserNurseByID(userID string) (entity.User, error) {
	if !isValidUUID(userID) {
		return entity.User{}, helper.CustomError{
			Code:    404,
			Message: errNurseNotFound.Error(),
		}
	}

	return u.userRepository.GetUserNurseByID(userID)
}
