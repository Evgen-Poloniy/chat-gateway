package v1

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine, handler *Handler) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/register", handler.RegisterUser)
		v1.POST("/create/chat/direct", handler.CreateDirectChat)
		v1.POST("/create/chat/group", handler.CreateGroupChat)
	}
}
