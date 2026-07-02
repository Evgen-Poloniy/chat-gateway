package router

import (
	"chat-gateway/internal/config"
	"chat-gateway/internal/middleware"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewRouter(logger *logrus.Logger, config *config.CORSConfig) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.HandleMethodNotAllowed = true

	router.Use(
		gin.Recovery(),
		middleware.CORS(config),
		middleware.SecureHeaders(),
		middleware.ErrorHandler(),
		middleware.Logger(logger),
	)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.NoRoute(func(c *gin.Context) {
		c.Error(fmt.Errorf("endpoint with URL %s not found", c.Request.URL.Path))
		c.Abort()
	})

	router.NoMethod(func(c *gin.Context) {
		c.Error(fmt.Errorf("method %s for URL %s not allowed", c.Request.Method, c.Request.URL.Path))
		c.Abort()
	})

	return router
}
