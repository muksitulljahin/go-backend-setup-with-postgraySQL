package imagecollection

import (
	"github.com/go-backend-setup-with-postgray/api/utils"
	"gorm.io/gorm"
)

// Standard content type constants for cross-service image management
const (
	ContentTypeProductItem = "product_item"
	ContentTypeUser        = "user"
	ContentTypeAdmin       = "admin"
	ContentTypeCategory    = "category"
	ContentTypeBanner      = "banner"
)

// BeforeCreate GORM hook to assign a 24-character hexadecimal unique identifier
func (img *ImageCollection) BeforeCreate(tx *gorm.DB) (err error) {
	if img.ID == "" {
		img.ID = utils.GenerateID()
	}
	return nil
}

// ToResponseDTO maps ImageCollection entity to ImageResponseDTO
func (img *ImageCollection) ToResponseDTO() ImageResponseDTO {
	return ImageResponseDTO{
		ID:          img.ID,
		ImageURL:    img.ImageURL,
		ContentID:   img.ContentID,
		ContentType: img.ContentType,
		IsPrimary:   img.IsPrimary,
		SortOrder:   img.SortOrder,
		CreatedAt:   img.CreatedAt,
		UpdatedAt:   img.UpdatedAt,
	}
}
