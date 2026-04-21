package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth(validKeys []string) gin.HandlerFunc {
	keySet := make(map[string]bool, len(validKeys))
	for _, k := range validKeys {
		keySet[strings.TrimSpace(k)] = true
	}
	return func(c *gin.Context) {
		if len(keySet) == 0 {
			c.Next()
			return
		}
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    40101,
				"message": "missing or invalid authorization header",
				"data":    nil,
			})
			c.Abort()
			return
		}
		key := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		if !keySet[key] {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    40101,
				"message": "invalid API key",
				"data":    nil,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
