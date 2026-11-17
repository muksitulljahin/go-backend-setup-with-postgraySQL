package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-backend-setup-with-postgray/pkg/jwt"
	"github.com/go-backend-setup-with-postgray/response"
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

// RequireRoles verifies if the authenticated user has one of the allowed roles
func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("userRole")
		if !exists {
			response.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		userRole, ok := roleVal.(string)
		if !ok {
			response.Forbidden(c, "Invalid user role format")
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if strings.EqualFold(userRole, allowed) {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "Access denied: insufficient permissions for this resource")
		c.Abort()
	}
}
