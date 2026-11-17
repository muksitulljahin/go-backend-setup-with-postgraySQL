package routes

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-backend-setup-with-postgray/api/middlewares"
	imagecollection "github.com/go-backend-setup-with-postgray/api/v1/imageCollection"
	"github.com/go-backend-setup-with-postgray/config"
	"github.com/go-backend-setup-with-postgray/response"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type swaggerHTMLModifier struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *swaggerHTMLModifier) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *swaggerHTMLModifier) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}

func SetupRouter(db *gorm.DB, cfg *config.AppConfig) *gin.Engine {
	r := gin.Default()

	// CORS Middleware with origin validation
	r.Use(middlewares.CORSMiddleware(cfg))

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

	// Middleware to inject multiple file selection script into Swagger UI
	swaggerGroup.Use(func(c *gin.Context) {
		blw := &swaggerHTMLModifier{body: bytes.NewBuffer(nil), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()

		contentType := blw.Header().Get("Content-Type")
		bodyBytes := blw.body.Bytes()
		if strings.Contains(contentType, "text/html") {
			injectedScript := `<script>
document.addEventListener("DOMContentLoaded", function() {
    var obs = new MutationObserver(function() {
        var files = document.querySelectorAll('input[type="file"]');
        for (var i = 0; i < files.length; i++) {
            if (!files[i].hasAttribute("multiple")) {
                files[i].setAttribute("multiple", "multiple");
                files[i].setAttribute("accept", "image/*");
            }
        }
    });
    obs.observe(document.body, { childList: true, subtree: true });
});
</script>
</body>`
			modified := strings.Replace(string(bodyBytes), "</body>", injectedScript, 1)
			blw.ResponseWriter.Header().Set("Content-Length", fmt.Sprintf("%d", len(modified)))
			blw.ResponseWriter.Write([]byte(modified))
		} else {
			blw.ResponseWriter.Write(bodyBytes)
		}
	})

	swaggerGroup.GET("/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.PersistAuthorization(true),
		ginSwagger.DocExpansion("none"),
	))

	// API v1 group (Add your new feature modules here!)
	apiV1 := r.Group("/api/v1")
	{
		imagecollection.RegisterImageRoutes(apiV1, db, cfg)
	}

	return r
}
