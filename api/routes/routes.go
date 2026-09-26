package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/muksitulljahin/go-backend-setup-with-postgraySQL/config"
	_ "github.com/muksitulljahin/go-backend-setup-with-postgraySQL/docs"
	"github.com/muksitulljahin/go-backend-setup-with-postgraySQL/response"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.AppConfig) *gin.Engine {
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

	// Root Route for fake message but still working
	r.GET("/", func(c *gin.Context) {
		response.Error(c, http.StatusInternalServerError, "Internal Server Error", nil)
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, "Health check passed", gin.H{
			"status": "UP",
		})
	})

	// Swagger API Documentation UI (Protected with Basic Auth)
	swaggerGroup := r.Group("/swagger")
	if cfg != nil && cfg.SwaggerUser != "" && cfg.SwaggerPassword != "" {
		swaggerGroup.Use(gin.BasicAuth(gin.Accounts{
			cfg.SwaggerUser: cfg.SwaggerPassword,
		}))
	}
	swaggerGroup.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 group (Add your new feature modules here!)
	apiV1 := r.Group("/api/v1")
	{
		_ = apiV1
		// Example: product.RegisterRoutes(apiV1, db)
	}

	return r
}
