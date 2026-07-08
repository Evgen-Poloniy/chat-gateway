package v1

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine, handler *Handler) {
	v1 := router.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/register", handler.RegisterUser)
	}

	protected := v1.Group("")
	{
		chats := protected.Group("/chats")
		{
			chats.POST("/create/direct", handler.CreateGroupChat)
			chats.POST("/create/group", handler.CreateGroupChat)
			chats.POST("/send", handler.SendMessage)
			chats.GET("/get/:user_id", handler.GetChatsByUserID)
		}

		users := protected.Group("/users")
		{
			users.GET("/search/:username", handler.GetUserDataByUsername)
			users.GET("/get/:user_id", handler.GetUserDataByUserID)
			users.PATCH("/update/:user_id", handler.UpdateUser)
		}
	}
}
