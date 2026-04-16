package middleware

import "github.com/gin-gonic/gin"

func SetupGinLoggerMiddleware() gin.HandlerFunc {
	logger := gin.Logger()
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}
		logger(c)
	}
}
