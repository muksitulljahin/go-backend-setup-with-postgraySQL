package imagecollection

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-backend-setup-with-postgray/response"
)

type ImageCollectionController struct {
	service *ImageCollectionService
}

func NewImageCollectionController(service *ImageCollectionService) *ImageCollectionController {
	return &ImageCollectionController{service: service}
}

func (ctrl *ImageCollectionController) UploadImage(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "File is required (key: 'file')", err.Error())
		return
	}

	contentID := c.PostForm("contentId")
	contentType := c.PostForm("contentType")
	isPrimaryStr := c.PostForm("isPrimary")
	isPrimary := isPrimaryStr == "true" || isPrimaryStr == "1"

	sortOrder := 0
	if sortOrderStr := c.PostForm("sortOrder"); sortOrderStr != "" {
		if val, parseErr := strconv.Atoi(sortOrderStr); parseErr == nil {
			sortOrder = val
		}
	}

	result, err := ctrl.service.UploadImage(c.Request.Context(), fileHeader, contentID, contentType, isPrimary, sortOrder)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Created(c, "Image uploaded to Cloudflare R2 and saved successfully", result)
}

func (ctrl *ImageCollectionController) UploadBulkImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.BadRequest(c, "Multipart form is required", err.Error())
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		if singleFile := form.File["file"]; len(singleFile) > 0 {
			files = singleFile
		} else {
			response.BadRequest(c, "At least one file is required (field: 'files')", nil)
			return
		}
	}

	contentID := c.PostForm("contentId")
	contentType := c.PostForm("contentType")

	results, err := ctrl.service.UploadBulkImages(c.Request.Context(), files, contentID, contentType)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Created(c, "Images uploaded to Cloudflare R2 and saved successfully", results)
}

func (ctrl *ImageCollectionController) UploadSingleImage(c *gin.Context) {
	contentID := c.Param("contentId")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "File is required (key: 'file')", err.Error())
		return
	}

	contentType := c.PostForm("contentType")

	result, err := ctrl.service.UploadAndUpsertSingleImage(c.Request.Context(), fileHeader, contentID, contentType)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, "Image uploaded to Cloudflare R2 and updated successfully", result)
}

func (ctrl *ImageCollectionController) CreateImage(c *gin.Context) {
	var dto CreateImageDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, "Invalid input data", err.Error())
		return
	}

	result, err := ctrl.service.CreateImage(dto)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Created(c, "Image added successfully", result)
}

func (ctrl *ImageCollectionController) UpdateImage(c *gin.Context) {
	id := c.Param("id")

	var dto UpdateImageDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, "Invalid input data", err.Error())
		return
	}

	result, err := ctrl.service.UpdateImage(id, dto)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, "Image updated successfully", result)
}

func (ctrl *ImageCollectionController) UpsertSingleImage(c *gin.Context) {
	contentID := c.Param("contentId")

	var dto UpsertSingleImageDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, "Invalid input data", err.Error())
		return
	}

	result, err := ctrl.service.UpsertSingleImage(contentID, dto)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, "Image saved successfully", result)
}

func (ctrl *ImageCollectionController) CreateBulkImages(c *gin.Context) {
	var dto CreateBulkImagesDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, "Invalid input data", err.Error())
		return
	}

	results, err := ctrl.service.CreateBulkImages(dto)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Created(c, "Images added successfully", results)
}

func (ctrl *ImageCollectionController) GetImagesByContentID(c *gin.Context) {
	contentID := c.Param("contentId")
	contentType := c.Query("contentType")

	results, err := ctrl.service.GetImagesByContentID(contentID, contentType)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, "Images retrieved successfully", results)
}

func (ctrl *ImageCollectionController) DeleteImage(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.service.DeleteImage(id); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, "Image deleted successfully", nil)
}
