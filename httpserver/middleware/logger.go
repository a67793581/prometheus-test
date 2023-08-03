package middleware

import (
	"fmt"
	"strings"
	"time"

	"prometheus-test/lib/logger"
	"prometheus-test/lib/util"

	"github.com/gin-gonic/gin"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()

		cost := time.Since(start) //

		logger.AccessInfo(c, fmt.Sprintf("%s|%s|%.3f|%s|%s|%s|%d|%s|%s|%d|%s|%s|%s|%s|%s|%s|%s|\n",
			start.Format("2006-01-02T15:04:05.000Z07:00"),
			clientIP(c),
			cost.Seconds(),
			"-", // optional
			"-", // optional
			"-", // optional
			c.Writer.Status(),
			"-", // optional
			c.Request.Host,
			c.Writer.Size(),
			path,
			c.Request.Method,
			util.GetRequestId(c),
			"-", // reserved
			"-", // reserved
			"-", // reserved
			"-", // reserved
		))
	}
}

func clientIP(ctx *gin.Context) string {
	ip := ctx.GetHeader("X-Forwarded-For")
	if index := strings.IndexByte(ip, ','); index >= 0 {
		ip = ip[0:index]
	}
	ip = strings.TrimSpace(ip)
	if len(ip) > 0 {
		return ip
	}
	ip = strings.TrimSpace(ctx.GetHeader("X-Real-Ip"))
	if len(ip) > 0 {
		return ip
	}
	return ctx.RemoteIP()
}
