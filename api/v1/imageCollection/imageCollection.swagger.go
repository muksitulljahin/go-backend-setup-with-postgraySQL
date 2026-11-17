package imagecollection

import "github.com/go-backend-setup-with-postgray/response"

// Compile-time check to ensure response package is used
var _ = response.APIResponse{}

// UploadImage godoc
// @Summary Upload Single Image to Cloudflare R2 and Save to DB
// @Description Upload a multipart file to Cloudflare R2 bucket, then store the resulting URL in ImageCollection
// @Tags Image Collection
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Image File to upload"
// @Param contentId formData string true "Content ID (e.g. ProductItem ID, User ID, Category ID)"
// @Param contentType formData string true "Content Type (product_item, user, admin, category)"
// @Param isPrimary formData bool false "Mark as Primary Image"
// @Param sortOrder formData int false "Display Sort Order"
// @Success 201 {object} response.APIResponse{data=ImageResponseDTO}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/images/upload [post]
func _UploadImageDoc() {}

// UploadBulkImages godoc
// @Summary Upload Multiple Images to Cloudflare R2 and Save to DB
// @Description Upload multiple multipart files to Cloudflare R2 bucket and store all resulting URLs in ImageCollection
// @Tags Image Collection
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param files formData file true "Multiple Image Files (field name: 'files')"
// @Param contentId formData string true "Content ID (e.g. ProductItem ID)"
// @Param contentType formData string true "Content Type (product_item, user, admin, category)"
// @Success 201 {object} response.APIResponse{data=[]ImageResponseDTO}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/images/upload/bulk [post]
func _UploadBulkImagesDoc() {}

// UploadSingleImage godoc
// @Summary Upload & Replace Single Image (Avatar/Banner)
// @Description Upload a single image file to Cloudflare R2 and upsert it for User, Admin, or Category
// @Tags Image Collection
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param contentId path string true "Content ID (e.g. User ID, Admin ID, Category ID)"
// @Param file formData file true "Avatar or Banner image file (field: 'file')"
// @Param contentType formData string true "Content Type (user, admin, category)"
// @Success 200 {object} response.APIResponse{data=ImageResponseDTO}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/images/upload/content/{contentId} [put]
func _UploadSingleImageDoc() {}

// CreateImage godoc
// @Summary Add Single Image to Content
// @Description Associate an image URL with any content ID (User, Admin, ProductItem, Category)
// @Tags Image Collection
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateImageDTO true "Create Image Payload"
// @Success 201 {object} response.APIResponse{data=ImageResponseDTO}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/images [post]
func _CreateImageDoc() {}

// UpdateImage godoc
// @Summary Update Image by ID
// @Description Update an existing image record (URL, sort order, primary status)
// @Tags Image Collection
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Image ID"
// @Param request body UpdateImageDTO true "Update Image Payload"
// @Success 200 {object} response.APIResponse{data=ImageResponseDTO}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/images/{id} [put]
func _UpdateImageDoc() {}

// UpsertSingleImage godoc
// @Summary Set or Update Single Image (Avatar/Banner)
// @Description Set or update single image for any content (User avatar, Admin avatar, Category banner)
// @Tags Image Collection
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param contentId path string true "Content ID (e.g. User ID, Admin ID, Category ID)"
// @Param request body UpsertSingleImageDTO true "Upsert Image Payload"
// @Success 200 {object} response.APIResponse{data=ImageResponseDTO}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/images/content/{contentId} [put]
func _UpsertSingleImageDoc() {}

// CreateBulkImages godoc
// @Summary Add Multiple Images to Content
// @Description Associate an array of image URLs with a content ID (Admin and SuperAdmin)
// @Tags Image Collection
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateBulkImagesDTO true "Bulk Images Payload"
// @Success 201 {object} response.APIResponse{data=[]ImageResponseDTO}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/images/bulk [post]
func _CreateBulkImagesDoc() {}

// GetImagesByContentID godoc
// @Summary Get Images by Content ID
// @Description Retrieve all image URLs associated with a content ID (supports ?contentType= filter)
// @Tags Image Collection
// @Produce json
// @Param contentId path string true "Content ID (e.g. User ID, Product Item ID)"
// @Param contentType query string false "Filter by Content Type (e.g. user, admin, product_item)"
// @Success 200 {object} response.APIResponse{data=[]ImageResponseDTO}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/images/content/{contentId} [get]
func _GetImagesByContentIDDoc() {}

// DeleteImage godoc
// @Summary Delete Image by ID
// @Description Delete a specific image entry by its ID (Admin and SuperAdmin)
// @Tags Image Collection
// @Produce json
// @Security BearerAuth
// @Param id path string true "Image ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/images/{id} [delete]
func _DeleteImageDoc() {}
