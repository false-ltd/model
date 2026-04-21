package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		reqID, _ := c.Get("request_id")

		if query != "" {
			path = path + "?" + query
		}

		log.Printf("[%s] %d %s %s %v", reqID, status, c.Request.Method, path, latency)
	}
}
