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
			chats.POST("/create/direct", handler.CreateDirectChatReq)
			chats.POST("/create/group", handler.CreateGroupChatReq)
			chats.POST("/send", handler.SendMessage)
		}

		users := protected.Group("/users")
		{
			users.GET("/get/:username", handler.GetUserDataByUsername)
			users.PATCH("/update/:id", handler.UpdateUser)
		}
	}
}
