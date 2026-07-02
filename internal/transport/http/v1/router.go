package v1

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/", nil)
	}
}
