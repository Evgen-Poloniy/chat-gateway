package router

import (
	"chat-gateway/internal/config"
	"chat-gateway/internal/dto"
	"chat-gateway/internal/middleware"
	"chat-gateway/internal/transport/ws"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// NewRouter initializes HTTP API router and connects middlewares on target routers.
func NewRouter(logger *logrus.Logger, config *config.CORSConfig) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.HandleMethodNotAllowed = true

	router.NoRoute(func(c *gin.Context) {
		statusCode := http.StatusNotFound
		code := "NOT_FOUND"
		message := fmt.Sprintf("endpoint with URL %s not found", c.Request.URL.Path)
		logAndResponseJSON(c, logger, statusCode, code, message)
	})

	router.NoMethod(func(c *gin.Context) {
		statusCode := http.StatusMethodNotAllowed
		code := "METHOD_NOT_ALLOWED"
		message := fmt.Sprintf("method %s for URL %s not allowed", c.Request.Method, c.Request.URL.Path)
		logAndResponseJSON(c, logger, statusCode, code, message)
	})

	router.Use(
		gin.Recovery(),
		middleware.CORS(config),
		middleware.SecureHeaders(),
	)

	wsRouter := router.Group("/")
	{
		wsRouter.GET("/ws", func(c *gin.Context) {
			ws.WebSocketUpgrade(c, logger)
		})
	}

	api := router.Group("/")
	api.Use(
		middleware.Logger(logger),
		middleware.ErrorHandler(),
	)

	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return router
}

// logAndResponseJSON is used in endpoints for logging and request request errors:
// URL not found and method not allowed.
func logAndResponseJSON(c *gin.Context, logger *logrus.Logger, statusCode int, code string, message string) {
	id := uuid.New().String()

	entry := logger.WithFields(map[string]interface{}{
		"id":     id,
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	})
	entry.Info("request start")

	start := time.Now()

	entry = logger.WithFields(map[string]interface{}{
		"id":          id,
		"method":      c.Request.Method,
		"path":        c.Request.URL.Path,
		"ip":          c.ClientIP(),
		"latency":     fmt.Sprintf("%v", time.Since(start)),
		"status_code": statusCode,
		"code":        code,
	})
	entry.Warn(message)

	c.AbortWithStatusJSON(statusCode, dto.ResponseError{
		Error: dto.Error{
			Code:    code,
			Message: message,
		},
	})
}
