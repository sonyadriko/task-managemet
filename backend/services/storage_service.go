package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type StorageService struct {
	client     *s3.Client
	bucketName string
	publicURL  string
}

func NewStorageService() (*StorageService, error) {
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("R2_BUCKET_NAME")
	publicURL := os.Getenv("R2_PUBLIC_URL")

	if accountID == "" || accessKeyID == "" || secretAccessKey == "" || bucketName == "" {
		return nil, fmt.Errorf("missing R2 configuration environment variables")
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &StorageService{
		client:     client,
		bucketName: bucketName,
		publicURL:  publicURL,
	}, nil
}

// Upload uploads a file to R2 and returns the storage key
func (s *StorageService) Upload(ctx context.Context, file io.Reader, originalFilename, mimeType string) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(originalFilename)
	storageKey := fmt.Sprintf("attachments/%s/%s%s", time.Now().Format("2006/01"), uuid.New().String(), ext)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(storageKey),
		Body:        file,
		ContentType: aws.String(mimeType),
	})
	if err != nil {
		return "", err
	}

	return storageKey, nil
}

// GetPresignedURL generates a presigned URL for downloading a file
func (s *StorageService) GetPresignedURL(ctx context.Context, storageKey string) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(storageKey),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", err
	}

	return request.URL, nil
}

// Delete removes a file from R2
func (s *StorageService) Delete(ctx context.Context, storageKey string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(storageKey),
	})
	return err
}

// GetPublicURL returns public URL if bucket is public, otherwise empty
func (s *StorageService) GetPublicURL(storageKey string) string {
	if s.publicURL != "" {
		return fmt.Sprintf("%s/%s", s.publicURL, storageKey)
	}
	return ""
}
