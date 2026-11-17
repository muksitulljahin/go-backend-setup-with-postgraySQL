package imagecollection

import "time"

// ImageCollection represents an image record associated with any parent content (User, Admin, ProductItem, Category, etc.)
type ImageCollection struct {
	ID          string    `gorm:"primaryKey;type:varchar(24)" json:"id" example:"670c5e7b8a1c9f2d3e4b5a6c"`
	ImageURL    string    `gorm:"type:text;not null" json:"imageUrl" example:"https://example.com/images/avatar.jpg"`
	ContentID   string    `gorm:"type:varchar(24);index;not null" json:"contentId" example:"65287f33a105f9cb2467d12b"`
	ContentType string    `gorm:"type:varchar(50);index;not null" json:"contentType" example:"user"`
	IsPrimary   bool      `gorm:"type:boolean;default:false" json:"isPrimary"`
	SortOrder   int       `gorm:"default:0" json:"sortOrder"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CreateImageDTO defines payload for creating a single image entry
type CreateImageDTO struct {
	ImageURL    string `json:"imageUrl" binding:"required" example:"https://example.com/images/avatar.jpg"`
	ContentID   string `json:"contentId" binding:"required" example:"65287f33a105f9cb2467d12b"`
	ContentType string `json:"contentType" binding:"required" example:"user"`
	IsPrimary   bool   `json:"isPrimary" example:"true"`
	SortOrder   int    `json:"sortOrder" example:"0"`
}

// UpdateImageDTO defines payload for updating an existing image record by ID
type UpdateImageDTO struct {
	ImageURL    string `json:"imageUrl" example:"https://example.com/images/updated-avatar.jpg"`
	ContentType string `json:"contentType" example:"user"`
	IsPrimary   *bool  `json:"isPrimary" example:"true"`
	SortOrder   *int   `json:"sortOrder" example:"1"`
}

// UpsertSingleImageDTO defines payload for setting or updating a single image (e.g. User or Admin avatar)
type UpsertSingleImageDTO struct {
	ImageURL    string `json:"imageUrl" binding:"required" example:"https://example.com/images/avatar.jpg"`
	ContentType string `json:"contentType" binding:"required" example:"user"`
}

// CreateBulkImagesDTO defines payload for adding multiple images to a content item
type CreateBulkImagesDTO struct {
	ContentID   string   `json:"contentId" binding:"required" example:"65287f33a105f9cb2467d12b"`
	ContentType string   `json:"contentType" binding:"required" example:"product_item"`
	ImageURLs   []string `json:"imageUrls" binding:"required,min=1" swaggertype:"array,string" example:"https://example.com/img1.jpg"`
}

// ImageResponseDTO defines sanitized response data for an image
type ImageResponseDTO struct {
	ID          string    `json:"id" example:"670c5e7b8a1c9f2d3e4b5a6c"`
	ImageURL    string    `json:"imageUrl" example:"https://example.com/images/avatar.jpg"`
	ContentID   string    `json:"contentId" example:"65287f33a105f9cb2467d12b"`
	ContentType string    `json:"contentType" example:"user"`
	IsPrimary   bool      `json:"isPrimary" example:"true"`
	SortOrder   int       `json:"sortOrder" example:"0"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
