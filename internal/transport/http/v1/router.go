package v1

import (
	"github.com/Evgen-Poloniy/chat-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine, handler *Handler) {
	v1 := router.Group("/api/v1")
	v1.Use(
		middleware.Logger(handler.logger),
		middleware.ErrorHandler(),
	)

	webhooks := v1.Group("/webhooks")
	{
		casdoor := webhooks.Group("/casdoor")
		casdoor.Use(middleware.APIKeyAuth(handler.config.IdPApiKey))
		{
			casdoor.POST("/users", handler.SingUpUser)
			casdoor.PUT("/users", handler.UpdateUser)
			casdoor.DELETE("/users", handler.DeleteUser)
		}
	}

	protected := v1.Group("")
	protected.Use(
		middleware.APIKeyAuth(handler.config.ApiKey),
		middleware.AuthMiddleware(handler.jwks, handler.config),
	)
	{
		chats := protected.Group("/chats")
		{
			chats.GET("/:user_id", handler.GetChatsByUserID)
			chats.POST("/direct", handler.CreateDirectChat)
			chats.POST("/group", handler.CreateGroupChat)
			chats.PATCH("/:chat_id", handler.UpdateGroupChat)
			chats.POST("/:chat_id/messages", handler.SendMessage)
		}

		users := protected.Group("/users")
		{
			users.GET("", handler.SearchUser)
			users.GET("/:user_id", handler.GetUserDataByUserID)
		}
	}
}
