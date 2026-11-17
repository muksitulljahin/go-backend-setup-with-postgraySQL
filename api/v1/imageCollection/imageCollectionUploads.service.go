package imagecollection

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"strings"

	r2bucket "github.com/go-backend-setup-with-postgray/api/cloude/r2Bucket"
	"github.com/go-backend-setup-with-postgray/response"
	"gorm.io/gorm"
)

type ImageCollectionService struct {
	db        *gorm.DB
	r2Service *r2bucket.R2Service
}

func NewImageCollectionService(db *gorm.DB, r2Service *r2bucket.R2Service) *ImageCollectionService {
	return &ImageCollectionService{
		db:        db,
		r2Service: r2Service,
	}
}

// UploadImage uploads an image file to Cloudflare R2 and saves its record into ImageCollection in DB
func (s *ImageCollectionService) UploadImage(ctx context.Context, fileHeader *multipart.FileHeader, contentID string, contentType string, isPrimary bool, sortOrder int) (*ImageResponseDTO, error) {
	if s.r2Service == nil {
		return nil, response.NewAppError(http.StatusInternalServerError, "Cloudflare R2 storage service is not initialized")
	}

	trimmedContentID := strings.TrimSpace(contentID)
	if trimmedContentID == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "Content ID is required")
	}

	trimmedContentType := strings.TrimSpace(contentType)
	if trimmedContentType == "" {
		trimmedContentType = "general"
	}

	uploadResult, err := s.r2Service.UploadFile(ctx, fileHeader, trimmedContentType)
	if err != nil {
		return nil, response.NewAppError(http.StatusInternalServerError, "Failed to upload image to Cloudflare R2", err)
	}

	return s.CreateImage(CreateImageDTO{
		ImageURL:    uploadResult.URL,
		ContentID:   trimmedContentID,
		ContentType: trimmedContentType,
		IsPrimary:   isPrimary,
		SortOrder:   sortOrder,
	})
}

// UploadBulkImages uploads multiple image files to Cloudflare R2 and saves all records into ImageCollection in DB
func (s *ImageCollectionService) UploadBulkImages(ctx context.Context, files []*multipart.FileHeader, contentID string, contentType string) ([]ImageResponseDTO, error) {
	if s.r2Service == nil {
		return nil, response.NewAppError(http.StatusInternalServerError, "Cloudflare R2 storage service is not initialized")
	}

	trimmedContentID := strings.TrimSpace(contentID)
	if trimmedContentID == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "Content ID is required")
	}

	trimmedContentType := strings.TrimSpace(contentType)
	if trimmedContentType == "" {
		trimmedContentType = "general"
	}

	if len(files) == 0 {
		return nil, response.NewAppError(http.StatusBadRequest, "At least one image file must be provided")
	}

	uploadResults, err := s.r2Service.UploadMultipleFiles(ctx, files, trimmedContentType)
	if err != nil {
		return nil, response.NewAppError(http.StatusInternalServerError, "Failed to upload images to Cloudflare R2", err)
	}

	imageURLs := make([]string, len(uploadResults))
	for i, res := range uploadResults {
		imageURLs[i] = res.URL
	}

	return s.CreateBulkImages(CreateBulkImagesDTO{
		ContentID:   trimmedContentID,
		ContentType: trimmedContentType,
		ImageURLs:   imageURLs,
	})
}

// UploadAndUpsertSingleImage uploads a single file to R2 and updates/sets single avatar for User, Admin, or Category
func (s *ImageCollectionService) UploadAndUpsertSingleImage(ctx context.Context, fileHeader *multipart.FileHeader, contentID string, contentType string) (*ImageResponseDTO, error) {
	if s.r2Service == nil {
		return nil, response.NewAppError(http.StatusInternalServerError, "Cloudflare R2 storage service is not initialized")
	}

	trimmedContentID := strings.TrimSpace(contentID)
	if trimmedContentID == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "Content ID is required")
	}

	trimmedContentType := strings.TrimSpace(contentType)
	if trimmedContentType == "" {
		trimmedContentType = "general"
	}

	uploadResult, err := s.r2Service.UploadFile(ctx, fileHeader, trimmedContentType)
	if err != nil {
		return nil, response.NewAppError(http.StatusInternalServerError, "Failed to upload image to Cloudflare R2", err)
	}

	return s.UpsertSingleImage(trimmedContentID, UpsertSingleImageDTO{
		ImageURL:    uploadResult.URL,
		ContentType: trimmedContentType,
	})
}

// UpsertSingleImage sets or updates a single image for an entity (e.g. User avatar, Admin avatar, Category banner)
// If an image already exists for the (contentID, contentType) pair, it updates it. Otherwise, it creates a new one.
func (s *ImageCollectionService) UpsertSingleImage(contentID string, dto UpsertSingleImageDTO) (*ImageResponseDTO, error) {
	trimmedContentID := strings.TrimSpace(contentID)
	if trimmedContentID == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "Content ID is required")
	}

	trimmedURL := strings.TrimSpace(dto.ImageURL)
	if trimmedURL == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "Image URL is required")
	}

	trimmedContentType := strings.TrimSpace(dto.ContentType)
	if trimmedContentType == "" {
		trimmedContentType = "general"
	}

	var existing ImageCollection
	err := s.db.Where("content_id = ? AND content_type = ?", trimmedContentID, trimmedContentType).First(&existing).Error
	if err == nil {
		// Update existing image record
		existing.ImageURL = trimmedURL
		existing.IsPrimary = true
		if err := s.db.Save(&existing).Error; err != nil {
			return nil, response.NewAppError(http.StatusInternalServerError, "Failed to update image record", err)
		}
		resp := existing.ToResponseDTO()
		return &resp, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, response.NewAppError(http.StatusInternalServerError, "Database error", err)
	}

	// Create new image record
	newImg := ImageCollection{
		ImageURL:    trimmedURL,
		ContentID:   trimmedContentID,
		ContentType: trimmedContentType,
		IsPrimary:   true,
		SortOrder:   0,
	}

	if err := s.db.Create(&newImg).Error; err != nil {
		return nil, response.NewAppError(http.StatusInternalServerError, "Failed to save image record", err)
	}

	resp := newImg.ToResponseDTO()
	return &resp, nil
}

// CreateBulkImages creates multiple image records for a content item
func (s *ImageCollectionService) CreateBulkImages(dto CreateBulkImagesDTO) ([]ImageResponseDTO, error) {
	trimmedContentID := strings.TrimSpace(dto.ContentID)
	if trimmedContentID == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "Content ID is required")
	}

	if len(dto.ImageURLs) == 0 {
		return nil, response.NewAppError(http.StatusBadRequest, "At least one image URL must be provided")
	}

	contentType := strings.TrimSpace(dto.ContentType)
	if contentType == "" {
		contentType = "general"
	}

	var images []ImageCollection
	for i, u := range dto.ImageURLs {
		trimmedU := strings.TrimSpace(u)
		if trimmedU == "" {
			continue
		}
		images = append(images, ImageCollection{
			ImageURL:    trimmedU,
			ContentID:   trimmedContentID,
			ContentType: contentType,
			IsPrimary:   i == 0,
			SortOrder:   i,
		})
	}

	if len(images) == 0 {
		return nil, response.NewAppError(http.StatusBadRequest, "No valid image URLs provided")
	}

	if err := s.db.Create(&images).Error; err != nil {
		return nil, response.NewAppError(http.StatusInternalServerError, "Failed to save images", err)
	}

	responses := make([]ImageResponseDTO, len(images))
	for i, img := range images {
		responses[i] = img.ToResponseDTO()
	}

	return responses, nil
}
