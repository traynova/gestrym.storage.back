package middleware

import (
	"fmt"
	"gestrym-storage/src/common/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupGinLoggerMiddleware() gin.HandlerFunc {
	log := utils.NewLogger()

	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery
		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		bodySize := c.Writer.Size()
		ginErrors := c.Errors.ByType(gin.ErrorTypePrivate).String()

		msg := fmt.Sprintf("[GIN] %d | %v | %s | %s %s | bytes:%d",
			status, latency, clientIP, method, path, bodySize)

		if ginErrors != "" {
			msg += fmt.Sprintf(" | errors: %s", ginErrors)
		}

		switch {
		case status >= 500:
			log.Error(msg)
		case status >= 400:
			log.Error("[HTTP_%d] %s %s desde %s — latencia: %v | bytes: %d%s",
				status, method, path, clientIP, latency, bodySize,
				func() string {
					if ginErrors != "" {
						return fmt.Sprintf(" | gin_errors: %s", ginErrors)
					}
					return ""
				}())
		case status >= 300:
			log.Warn(msg)
		default:
			log.Info(msg)
		}
	}
}
