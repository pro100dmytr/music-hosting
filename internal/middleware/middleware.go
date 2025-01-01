package middleware

import (
	"github.com/gin-gonic/gin"
	"music-hosting/internal/utils/jwtutils"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			c.Abort()
			return
		}

		tokenString = tokenString[len("Bearer "):]

		userID, err := jwtutils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
