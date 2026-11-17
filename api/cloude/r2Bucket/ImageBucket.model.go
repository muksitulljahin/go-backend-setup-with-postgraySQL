package r2bucket

import "time"

// UploadResult represents metadata of a file successfully uploaded to Cloudflare R2
type UploadResult struct {
	Key         string    `json:"key" example:"uploads/product_item/670c5e7b8a1c9f2d3e4b5a6c_image.jpg"`
	URL         string    `json:"url" example:"https://your-image-domain.com/uploads/product_item/670c5e7b8a1c9f2d3e4b5a6c_image.jpg"`
	Filename    string    `json:"filename" example:"image.jpg"`
	Size        int64     `json:"size" example:"204800"`
	ContentType string    `json:"contentType" example:"image/jpeg"`
	UploadedAt  time.Time `json:"uploadedAt"`
}
