package imagecollection

import (
	"net/http"
	"strings"

	"github.com/go-backend-setup-with-postgray/response"
)

// CreateImage creates a single image record associated with any content ID (User, Admin, ProductItem, etc.)
func (s *ImageCollectionService) CreateImage(dto CreateImageDTO) (*ImageResponseDTO, error) {
	trimmedURL := strings.TrimSpace(dto.ImageURL)
	if trimmedURL == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "Image URL is required")
	}

	trimmedContentID := strings.TrimSpace(dto.ContentID)
	if trimmedContentID == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "Content ID is required")
	}

	trimmedContentType := strings.TrimSpace(dto.ContentType)
	if trimmedContentType == "" {
		trimmedContentType = "general"
	}

	img := ImageCollection{
		ImageURL:    trimmedURL,
		ContentID:   trimmedContentID,
		ContentType: trimmedContentType,
		IsPrimary:   dto.IsPrimary,
		SortOrder:   dto.SortOrder,
	}

	if err := s.db.Create(&img).Error; err != nil {
		return nil, response.NewAppError(http.StatusInternalServerError, "Failed to create image record", err)
	}

	resp := img.ToResponseDTO()
	return &resp, nil
}
