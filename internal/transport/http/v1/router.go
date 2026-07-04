package v1

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine, handler *Handler) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/register", handler.RegisterUser)
	}
}
