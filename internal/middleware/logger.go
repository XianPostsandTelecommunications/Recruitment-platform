package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"lab-recruitment-platform/pkg/logger"
)

// RequestLoggerMiddleware 请求日志中间件
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		status := c.Writer.Status()
		logger.Infof("请求日志 | %s | %s | %d | %v | %s | %s",
			c.Request.Method,
			c.Request.URL.Path,
			status,
			latency,
			c.ClientIP(),
			c.Request.UserAgent(),
		)
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Errorf("请求错误 | %s | %v",
					c.Request.URL.Path,
					err.Error(),
				)
			}
		}
	}
} 