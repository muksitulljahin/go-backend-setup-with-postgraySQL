package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-backend-setup-with-postgray/config"
)

// CORSMiddleware strictly validates and sets CORS headers based on configured frontend URLs.
func CORSMiddleware(cfg *config.AppConfig) gin.HandlerFunc {
	allowedOrigins := make(map[string]bool)

	// Default frontend and backend URLs
	defaults := []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://localhost:3002",
		"http://localhost:3003",
		"http://localhost:8080",
		"http://192.168.0.110:8080",
	}
	for _, o := range defaults {
		allowedOrigins[strings.ToLower(strings.TrimSpace(o))] = true
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// If no Origin header (e.g. mobile app, curl, Postman), proceed normally
		if origin == "" {
			c.Next()
			return
		}

		cleanOrigin := strings.ToLower(strings.TrimSpace(origin))

		// Strict origin matching
		if allowedOrigins[cleanOrigin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, x-refresh-token")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Authorization, x-refresh-token")
			// c.Writer.Header().Set("Access-Control-Max-Age", "86400")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}

			c.Next()
			return
		}

		// If origin is not recognized/allowed
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Cors Issue!!! Origin is not allowed by CORS policy",
		})
	}
}
