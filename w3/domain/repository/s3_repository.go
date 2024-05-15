package repository

import (
	"halosuster/internal/database"
	"mime/multipart"
)

type s3Repository struct {
	s3Writer database.IS3Writer
}

func NewS3Repository(s database.IS3Writer) *s3Repository {
	return &s3Repository{
		s3Writer: s,
	}
}

type IS3Repository interface {
	UploadImage(file *multipart.FileHeader) (string, error)
}

func (s *s3Repository) UploadImage(file *multipart.FileHeader) (string, error) {
	return s.s3Writer.UploadImage(file)
}
