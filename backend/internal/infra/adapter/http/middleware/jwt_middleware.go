package middleware

import (
	"net/http"
	"strings"
	"task-management/internal/applications/ports/services"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware(jwtService services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")

		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or malformed token"})
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("user_id", claims.ID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
