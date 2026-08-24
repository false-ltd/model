package middleware

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS configures cross-origin access. With wildcard origins credentials
// are disabled (the spec forbids combining Allow-Credentials with "*");
// credentials are only enabled for explicitly listed origins.
func CORS(allowedOrigins string) gin.HandlerFunc {
	config := cors.Config{
		AllowMethods:  []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset"},
	}
	if allowedOrigins == "" || allowedOrigins == "*" {
		config.AllowAllOrigins = true
	} else {
		config.AllowOrigins = parseOrigins(allowedOrigins)
		config.AllowCredentials = true
	}
	return cors.New(config)
}

func parseOrigins(s string) []string {
	parts := strings.Split(s, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			origins = append(origins, p)
		}
	}
	return origins
}
