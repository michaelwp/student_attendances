package config

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
}

func NewS3Config() *S3Config {
	return &S3Config{
		Region:          os.Getenv("AWS_S3_REGION"),
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		BucketName:      os.Getenv("AWS_S3_BUCKET"),
	}
}

func (s *S3Config) NewS3Client() (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(s.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			s.AccessKeyID,
			s.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}

func (s *S3Config) UploadFile(client *s3.Client, key string, body []byte) error {
	_, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.BucketName,
		Key:    &key,
		Body:   bytes.NewReader(body),
	})
	return err
}

// GetObjectURL returns a non-signed public URL for the S3 object
func (s *S3Config) GetObjectURL(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s/",
		s.BucketName, s.Region, key)
}

// GetSignedURL generates a presigned URL for the S3 object with specified expiration time
func (s *S3Config) GetSignedURL(client *s3.Client, key string, expires time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(client)
	request, err := presignClient.PresignGetObject(context.TODO(),
		&s3.GetObjectInput{
			Bucket: &s.BucketName,
			Key:    &key,
		},
		func(opts *s3.PresignOptions) {
			opts.Expires = expires
		},
	)
	if err != nil {
		return "", err
	}
	return request.URL, nil
}
