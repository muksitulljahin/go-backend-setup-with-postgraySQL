package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/muksitulljahin/go-backend-setup-with-postgraySQL/pkg/jwt"
	"github.com/muksitulljahin/go-backend-setup-with-postgraySQL/response"
)

// JWTAuthMiddleware verifies Bearer tokens on protected routes
func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Invalid authorization format. Expected: Bearer <token>")
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := jwt.ValidateToken(tokenStr, secret)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired JWT token")
			c.Abort()
			return
		}

		// Store user info in Gin context for downstream handlers
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}
