package imagecollection

import (
	"log"

	"github.com/gin-gonic/gin"
	r2bucket "github.com/go-backend-setup-with-postgray/api/cloude/r2Bucket"
	"github.com/go-backend-setup-with-postgray/api/enum"
	"github.com/go-backend-setup-with-postgray/api/middlewares"
	"github.com/go-backend-setup-with-postgray/config"

	"gorm.io/gorm"
)

// RegisterImageRoutes registers all image collection endpoints.
// NOTE: AutoMigrate is deliberately omitted as per project instruction.
func RegisterImageRoutes(router *gin.RouterGroup, db *gorm.DB, cfg *config.AppConfig) {
	r2Service, err := r2bucket.NewR2Service(cfg)
	if err != nil {
		log.Printf("⚠️ Cloudflare R2 bucket service warning: %v\n", err)
	}

	service := NewImageCollectionService(db, r2Service)
	controller := NewImageCollectionController(service)

	jwtSecret := "default_jwt_secret"
	if cfg != nil && cfg.JWTSecret != "" {
		jwtSecret = cfg.JWTSecret
	}

	imagesGroup := router.Group("/images")
	{
		// Public route to view images by content ID (with optional ?contentType=)
		imagesGroup.GET("/content/:contentId", controller.GetImagesByContentID)

		// Protected management routes (Admin and SuperAdmin)
		protected := imagesGroup.Group("")
		protected.Use(
			middlewares.JWTAuthMiddleware(jwtSecret),
			middlewares.RequireRoles(
				string(enum.SuperAdmin),
				string(enum.Admin),
			),
		)
		{
			// JSON URL endpoints
			protected.POST("", controller.CreateImage)
			protected.POST("/bulk", controller.CreateBulkImages)
			protected.PUT("/:id", controller.UpdateImage)
			protected.PUT("/content/:contentId", controller.UpsertSingleImage)
			protected.DELETE("/:id", controller.DeleteImage)

			// Cloudflare R2 multipart file upload endpoints
			protected.POST("/upload", controller.UploadImage)
			protected.POST("/upload/bulk", controller.UploadBulkImages)
			protected.PUT("/upload/content/:contentId", controller.UploadSingleImage)
		}
	}
}
