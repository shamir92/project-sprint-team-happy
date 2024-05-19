package configuration

import (
	"os"
)

type s3Configuration struct {
	s3ID         string
	s3SecretKey  string
	s3BucketName string
	s3Region     string
}

func NewS3Configuration() *s3Configuration {
	return &s3Configuration{
		s3ID:         os.Getenv("AWS_ACCESS_KEY_ID"),
		s3SecretKey:  os.Getenv("AWS_SECRET_ACCESS_KEY"),
		s3BucketName: os.Getenv("AWS_S3_BUCKET_NAME"),
		s3Region:     os.Getenv("AWS_REGION"),
	}
}

type IS3Configuration interface {
	GetS3ID() string
	GetS3SecretKey() string
	GetS3BucketName() string
	GetS3Region() string
}

func (c *s3Configuration) GetS3ID() string {
	return c.s3ID
}

func (c *s3Configuration) GetS3SecretKey() string {
	return c.s3SecretKey
}

func (c *s3Configuration) GetS3BucketName() string {
	return c.s3BucketName
}

func (c *s3Configuration) GetS3Region() string {
	return c.s3Region
}
