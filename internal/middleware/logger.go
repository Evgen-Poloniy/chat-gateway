package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Logger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statusCode int

		id := uuid.NewString()

		entry := logger.WithFields(map[string]interface{}{
			"id":     id,
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"ip":     c.ClientIP(),
		})

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if val, exists := c.Get("status_code"); exists {
				if v, ok := val.(int); ok {
					statusCode = v
					entry = entry.WithField("status_code", statusCode)
				}
			}
			if val, exists := c.Get("code"); exists {
				if v, ok := val.(string); ok {
					entry = entry.WithField("code", v)
				}
			}

			if statusCode >= 500 {
				entry.Error(err.Error())
			} else {
				entry.Warn(err.Error())
			}

			return
		}

		switch {
		case statusCode >= 500:
			entry.Error("request start with error")
			return
		case statusCode >= 400:
			entry.Warn("request start with warning")
			return
		default:
			entry.Info("request start")
		}

		start := time.Now()

		c.Next()

		latency := time.Since(start)
		statusCode = c.Writer.Status()

		entry = logger.WithFields(map[string]interface{}{
			"id":          id,
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"ip":          c.ClientIP(),
			"latency":     fmt.Sprintf("%v", latency),
			"status_code": statusCode,
		})

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if val, exists := c.Get("status_code"); exists {
				if v, ok := val.(int); ok {
					statusCode = v
					entry = entry.WithField("status_code", statusCode)
				}
			}
			if val, exists := c.Get("code"); exists {
				if v, ok := val.(string); ok {
					entry = entry.WithField("code", v)
				}
			}

			if statusCode >= 500 {
				entry.Error(err.Error())
			} else {
				entry.Warn(err.Error())
			}

			return
		}

		switch {
		case statusCode >= 500:
			entry.Error("request completed with error")
		case statusCode >= 400:
			entry.Warn("request completed with warning")
		default:
			entry.Info("request completed successfully")
		}
	}
}
