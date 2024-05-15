package database

import (
	"context"
	"fmt"
	"halosuster/configuration"
	"halosuster/internal/helper"
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type s3Writer struct {
	s3Client *s3.Client
	s3Config configuration.IS3Configuration
}

type IS3Writer interface {
	UploadImage(fileHeader *multipart.FileHeader) (string, error)
}

func NewS3Writer(configS3 configuration.IS3Configuration) (*s3Writer, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(configS3.GetS3Region()))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create an S3 client
	s3Client := s3.NewFromConfig(cfg)
	return &s3Writer{
		s3Client: s3Client,
		s3Config: configS3,
	}, nil
}

func (w *s3Writer) UploadImage(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()

	if err != nil {
		return "", helper.CustomError{
			Message: err.Error(),
			Code:    500,
		}
	}
	defer file.Close()
	fileName := uuid.New().String()
	bucket := w.s3Config.GetS3BucketName()

	_, err = w.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(fileName),
		Body:        file,
		ACL:         types.ObjectCannedACLPublicRead, // Adjust permissions as needed
		ContentType: aws.String("image/jpeg"),        // Set content type
	})

	if err != nil {
		log.Fatalf("failed to upload file, %v", err)
		return "", helper.CustomError{
			Message: err.Error(),
			Code:    500,
		}
	}
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, w.s3Config.GetS3Region(), fileName)
	return url, nil

}
