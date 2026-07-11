package v1

import (
	"github.com/Evgen-Poloniy/chat-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine, handler *Handler, apiKeyHash []byte) {
	v1 := router.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/register", handler.RegisterUser)
	}

	protected := v1.Group("")
	protected.Use(middleware.APIKeyAuth(apiKeyHash))
	{
		chats := protected.Group("/chats")
		{
			chats.GET("/search/:user_id", handler.GetChatsByUserID)
			chats.POST("/create/direct", handler.CreateDirectChat)
			chats.POST("/create/group", handler.CreateGroupChat)
			chats.PATCH("/update/:chat_id", handler.UpdateGroupChat)
			chats.POST("/send", handler.SendMessage)
		}

		users := protected.Group("/users")
		{
			users.GET("/search/:username", handler.GetUserDataByUsername)
			users.GET("/search/:user_id", handler.GetUserDataByUserID)
			users.PATCH("/update/:user_id", handler.UpdateUser)
		}
	}
}
