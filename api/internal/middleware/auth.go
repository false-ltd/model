package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// APIKeyAuth validates Bearer tokens against the configured key set.
// allowPublic permits unauthenticated triggers (bounded by cooldown, rate
// limit and the sync mutex) — used for the footer refresh button in
// deployments that opt in. Otherwise fail closed: when no keys are
// configured the endpoint stays locked so a misconfigured deployment
// cannot expose the sync endpoint to everyone.
func APIKeyAuth(validKeys []string, allowPublic bool) gin.HandlerFunc {
	keys := make([]string, 0, len(validKeys))
	for _, k := range validKeys {
		if k = strings.TrimSpace(k); k != "" {
			keys = append(keys, k)
		}
	}
	return func(c *gin.Context) {
		if allowPublic {
			c.Next()
			return
		}
		if len(keys) == 0 {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"code":    50301,
				"message": "server is not configured with API keys; sync endpoint disabled",
				"data":    nil,
			})
			c.Abort()
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
		if !keyMatches(keys, key) {
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

// keyMatches compares the candidate against every configured key in
// constant time per comparison.
func keyMatches(keys []string, candidate string) bool {
	match := 0
	for _, k := range keys {
		if subtle.ConstantTimeCompare([]byte(k), []byte(candidate)) == 1 {
			match = 1
		}
	}
	return match == 1
}
