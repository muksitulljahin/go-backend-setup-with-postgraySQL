package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/muksitulljahin/go-backend-setup-with-postgraySQL/docs"
	"github.com/muksitulljahin/go-backend-setup-with-postgraySQL/internal/modules/user"
	"github.com/muksitulljahin/go-backend-setup-with-postgraySQL/pkg/response"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Root Route
	r.GET("/", func(c *gin.Context) {
		response.Success(c, "backend running successfully", gin.H{
			"status":  "OK",
			"version": "1.0.0",
		})
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, "Health check passed", gin.H{
			"uptime": "active",
		})
	})

	// Swagger API Documentation UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 group
	apiV1 := r.Group("/api/v1")
	{
		user.RegisterRoutes(apiV1, db)
	}

	return r
}
