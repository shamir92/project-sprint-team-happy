package usecase

import (
	"halosuster/domain/repository"
	"halosuster/internal/helper"
	"mime/multipart"
)

// import "halosuster/domain/repository"

type imageUsecase struct {
	s3Repository repository.IS3Repository
}

type IImageUsecase interface {
	UploadImage(file *multipart.FileHeader) (string, error)
}

func NewImageUsecase(s3Repository repository.IS3Repository) *imageUsecase {
	return &imageUsecase{
		s3Repository: s3Repository,
	}
}

func (u *imageUsecase) UploadImage(file *multipart.FileHeader) (string, error) {

	url, err := u.s3Repository.UploadImage(file)
	if err != nil {
		return "", helper.CustomError{
			Message: err.Error(),
			Code:    500,
		}
	}
	return url, nil
}
