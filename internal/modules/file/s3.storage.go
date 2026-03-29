package file

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Storage struct {
	client *s3.Client
	bucket string
}

func NewS3Storage(client *s3.Client, bucket string) StorageProvider {
	return &s3Storage{
		client: client,
		bucket: bucket,
	}
}

func (s *s3Storage) Save(content io.Reader, objectKey string, contentType string) (string, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectKey),
		Body:   content,
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	_, err := s.client.PutObject(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("error al subir a S3: %w", err)
	}

	return objectKey, err
}

func (s *s3Storage) Get(storageKey string) (io.ReadCloser, error) {
	output, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(storageKey),
	})
	if err != nil {
		return nil, fmt.Errorf("error al descargar de S3: %w", err)
	}

	return output.Body, nil
}
