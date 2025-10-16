package middleware

import (
	"net/http"
	"strings"
	"task-management/internal/applications/ports/services"
	"task-management/internal/domain"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware(jwtService services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"error":   "missing or malformed token",
			})
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"error":   "empty authorization token",
			})
			return
		}

		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"error":   err.Error(),
			})
			return
		}

		// Set user info ke context
		c.Set("user", claims)
		c.Next()
	}
}

func GetUserClaims(c *gin.Context) (*domain.JWTClaims, bool) {
	claims, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	userClaims, ok := claims.(*domain.JWTClaims)
	if !ok {
		return nil, false
	}

	return userClaims, true
}
