package imagecollection

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-backend-setup-with-postgray/response"
	"gorm.io/gorm"
)

// UpdateImage updates an existing image record by its primary key ID
func (s *ImageCollectionService) UpdateImage(id string, dto UpdateImageDTO) (*ImageResponseDTO, error) {
	trimmedID := strings.TrimSpace(id)
	if trimmedID == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "Image ID is required")
	}

	var img ImageCollection
	if err := s.db.First(&img, "id = ?", trimmedID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewAppError(http.StatusNotFound, "Image not found")
		}
		return nil, response.NewAppError(http.StatusInternalServerError, "Database error", err)
	}

	if dto.ImageURL != "" {
		img.ImageURL = strings.TrimSpace(dto.ImageURL)
	}

	if dto.ContentType != "" {
		img.ContentType = strings.TrimSpace(dto.ContentType)
	}

	if dto.IsPrimary != nil {
		img.IsPrimary = *dto.IsPrimary
	}

	if dto.SortOrder != nil {
		img.SortOrder = *dto.SortOrder
	}

	if err := s.db.Save(&img).Error; err != nil {
		return nil, response.NewAppError(http.StatusInternalServerError, "Failed to update image", err)
	}

	resp := img.ToResponseDTO()
	return &resp, nil
}
