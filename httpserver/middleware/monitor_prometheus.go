package middleware

import (
	"prometheus-test/metrics"
	"time"

	"github.com/gin-gonic/gin"
)

func MonitorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		// Process request
		c.Next()
		interval := time.Since(start).Milliseconds()
		statusCode := c.Writer.Status()

		metrics.UpdateInterfaceQPS(path, statusCode, 1)
		metrics.UpdateInterface(path, statusCode, interval)
	}
}
