package imagecollection

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-backend-setup-with-postgray/response"
	"gorm.io/gorm"
)

// GetImagesByContentID retrieves all images associated with a content ID (optionally filtered by contentType)
func (s *ImageCollectionService) GetImagesByContentID(contentID string, contentType ...string) ([]ImageResponseDTO, error) {
	trimmedContentID := strings.TrimSpace(contentID)
	if trimmedContentID == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "Content ID is required")
	}

	query := s.db.Where("content_id = ?", trimmedContentID)
	if len(contentType) > 0 && strings.TrimSpace(contentType[0]) != "" {
		query = query.Where("content_type = ?", strings.TrimSpace(contentType[0]))
	}

	var images []ImageCollection
	if err := query.Order("sort_order ASC, created_at ASC").Find(&images).Error; err != nil {
		return nil, response.NewAppError(http.StatusInternalServerError, "Failed to fetch images", err)
	}

	results := make([]ImageResponseDTO, len(images))
	for i, img := range images {
		results[i] = img.ToResponseDTO()
	}

	return results, nil
}

// SyncImages synchronizes the full image list for a content item (replaces previous images with new ones)
func (s *ImageCollectionService) SyncImages(contentID string, contentType string, imageURLs []string) ([]ImageResponseDTO, error) {
	trimmedContentID := strings.TrimSpace(contentID)
	if trimmedContentID == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "Content ID is required")
	}

	trimmedContentType := strings.TrimSpace(contentType)
	if trimmedContentType == "" {
		trimmedContentType = "general"
	}

	// Delete previous images for this content
	_ = s.db.Where("content_id = ? AND content_type = ?", trimmedContentID, trimmedContentType).Delete(&ImageCollection{}).Error

	if len(imageURLs) == 0 {
		return []ImageResponseDTO{}, nil
	}

	var images []ImageCollection
	for i, u := range imageURLs {
		trimmedU := strings.TrimSpace(u)
		if trimmedU == "" {
			continue
		}
		images = append(images, ImageCollection{
			ImageURL:    trimmedU,
			ContentID:   trimmedContentID,
			ContentType: trimmedContentType,
			IsPrimary:   i == 0,
			SortOrder:   i,
		})
	}

	if len(images) > 0 {
		if err := s.db.Create(&images).Error; err != nil {
			return nil, response.NewAppError(http.StatusInternalServerError, "Failed to save synchronized images", err)
		}
	}

	responses := make([]ImageResponseDTO, len(images))
	for i, img := range images {
		responses[i] = img.ToResponseDTO()
	}

	return responses, nil
}

// GetSingleImage retrieves the primary/first image for a content item (useful for avatar/thumbnail)
func (s *ImageCollectionService) GetSingleImage(contentID string, contentType string) (*ImageResponseDTO, error) {
	trimmedContentID := strings.TrimSpace(contentID)
	if trimmedContentID == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "Content ID is required")
	}

	var img ImageCollection
	query := s.db.Where("content_id = ?", trimmedContentID)
	if trimmedType := strings.TrimSpace(contentType); trimmedType != "" {
		query = query.Where("content_type = ?", trimmedType)
	}

	if err := query.Order("is_primary DESC, sort_order ASC, created_at ASC").First(&img).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No image found is not necessarily an application error
		}
		return nil, response.NewAppError(http.StatusInternalServerError, "Database error", err)
	}

	resp := img.ToResponseDTO()
	return &resp, nil
}
