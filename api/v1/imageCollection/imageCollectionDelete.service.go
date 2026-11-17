package imagecollection

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-backend-setup-with-postgray/response"
	"gorm.io/gorm"
)

// DeleteImage deletes a single image by ID
func (s *ImageCollectionService) DeleteImage(id string) error {
	trimmedID := strings.TrimSpace(id)
	if trimmedID == "" {
		return response.NewAppError(http.StatusBadRequest, "Image ID is required")
	}

	var img ImageCollection
	if err := s.db.First(&img, "id = ?", trimmedID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewAppError(http.StatusNotFound, "Image not found")
		}
		return response.NewAppError(http.StatusInternalServerError, "Database error", err)
	}

	if s.r2Service != nil && img.ImageURL != "" {
		_ = s.r2Service.DeleteFile(context.Background(), img.ImageURL)
	}

	if err := s.db.Delete(&img).Error; err != nil {
		return response.NewAppError(http.StatusInternalServerError, "Failed to delete image", err)
	}

	return nil
}

// DeleteImagesByContentID deletes all images for a given content item
func (s *ImageCollectionService) DeleteImagesByContentID(contentID string, contentType ...string) error {
	trimmedContentID := strings.TrimSpace(contentID)
	if trimmedContentID == "" {
		return response.NewAppError(http.StatusBadRequest, "Content ID is required")
	}

	query := s.db.Where("content_id = ?", trimmedContentID)
	if len(contentType) > 0 && strings.TrimSpace(contentType[0]) != "" {
		query = query.Where("content_type = ?", strings.TrimSpace(contentType[0]))
	}

	if err := query.Delete(&ImageCollection{}).Error; err != nil {
		return response.NewAppError(http.StatusInternalServerError, "Failed to delete content images", err)
	}

	return nil
}
