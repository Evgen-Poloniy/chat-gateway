package ws

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine, handler *Handler) {
	wsRouter := router.Group("/ws")
	{
		wsRouter.GET("", handler.WebSocketUpgrade)
	}
}
