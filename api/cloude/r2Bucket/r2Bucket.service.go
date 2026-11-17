package r2bucket

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-backend-setup-with-postgray/api/utils"
	"github.com/go-backend-setup-with-postgray/config"
)

type R2Service struct {
	client     *s3.Client
	bucketName string
	publicURL  string
	accountID  string
}

// NewR2Service initializes Cloudflare R2 S3-compatible storage service
func NewR2Service(cfg *config.AppConfig) (*R2Service, error) {
	if cfg == nil {
		return nil, errors.New("app config is required for R2 initialization")
	}

	accountID := strings.TrimSpace(cfg.R2AccountID)
	accessKey := strings.TrimSpace(cfg.R2AccessKeyID)
	secretKey := strings.TrimSpace(cfg.R2SecretAccessKey)
	bucketName := strings.TrimSpace(cfg.R2BucketName)

	if accountID == "" || accessKey == "" || secretKey == "" || bucketName == "" {
		return nil, errors.New("missing Cloudflare R2 credentials in environment variables (R2_ACCOUNT_ID, R2_ACCESS_KEY_ID, R2_SECRET_ACCESS_KEY, R2_BUCKET_NAME)")
	}

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)

	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		awsconfig.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config for R2: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return &R2Service{
		client:     s3Client,
		bucketName: bucketName,
		publicURL:  strings.TrimRight(strings.TrimSpace(cfg.R2PublicURL), "/"),
		accountID:  accountID,
	}, nil
}

// UploadFile uploads a multipart file to Cloudflare R2 and returns its public URL & metadata
func (s *R2Service) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (*UploadResult, error) {
	if fileHeader == nil {
		return nil, errors.New("file is required")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	cleanFolder := strings.Trim(strings.TrimSpace(folder), "/")
	if cleanFolder == "" {
		cleanFolder = "general"
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		ext = ".jpg"
	}

	uniqueID := utils.GenerateID()
	key := fmt.Sprintf("uploads/%s/%s%s", cleanFolder, uniqueID, ext)

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType([]byte{})
		if ext == ".png" {
			contentType = "image/png"
		} else if ext == ".webp" {
			contentType = "image/webp"
		} else if ext == ".gif" {
			contentType = "image/gif"
		} else {
			contentType = "image/jpeg"
		}
	}

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to Cloudflare R2: %w", err)
	}

	var publicURL string
	if s.publicURL != "" {
		publicURL = fmt.Sprintf("%s/%s", s.publicURL, key)
	} else {
		publicURL = fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s/%s", s.accountID, s.bucketName, key)
	}

	return &UploadResult{
		Key:         key,
		URL:         publicURL,
		Filename:    fileHeader.Filename,
		Size:        fileHeader.Size,
		ContentType: contentType,
		UploadedAt:  time.Now(),
	}, nil
}

// UploadMultipleFiles uploads a slice of multipart files to Cloudflare R2
func (s *R2Service) UploadMultipleFiles(ctx context.Context, files []*multipart.FileHeader, folder string) ([]UploadResult, error) {
	var results []UploadResult

	for _, fileHeader := range files {
		res, err := s.UploadFile(ctx, fileHeader, folder)
		if err != nil {
			return nil, err
		}
		results = append(results, *res)
	}

	return results, nil
}

// DeleteFile deletes a file from Cloudflare R2 using its full URL or key
func (s *R2Service) DeleteFile(ctx context.Context, fileKeyOrURL string) error {
	trimmed := strings.TrimSpace(fileKeyOrURL)
	if trimmed == "" {
		return nil
	}

	key := trimmed
	if s.publicURL != "" && strings.HasPrefix(key, s.publicURL) {
		key = strings.TrimPrefix(key, s.publicURL)
		key = strings.TrimPrefix(key, "/")
	} else if strings.Contains(key, s.bucketName+"/") {
		parts := strings.Split(key, s.bucketName+"/")
		if len(parts) == 2 {
			key = parts[1]
		}
	}

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from Cloudflare R2: %w", err)
	}

	return nil
}
