package v1

import (
	"chat-gateway/internal/dto"
	"chat-gateway/internal/entity"
	errs "chat-gateway/pkg/errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterUser(c *gin.Context) {
	var req dto.RegisterUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&errs.AppError{
			StatusCode: http.StatusBadRequest,
			Code:       "BAD_REQUEST",
			Message:    err.Error(),
		})
		return
	}

	user := &entity.User{
		Username:  req.Username,
		Email:     &req.Email,
		FirstName: &req.FirstName,
		LastName:  &req.LastName,
		BirthDate: &req.BirthDate,
	}

	if err := h.messenger.CreateUser(c.Request.Context(), user); err != nil {
		c.Error(&errs.AppError{
			StatusCode: http.StatusInternalServerError,
			Code:       "INTERNAL_SERVER_ERROR",
			Message:    "error when create user",
		})
	}

	userData := &dto.UserData{
		ID:        user.ID,
		Username:  user.Username,
		Email:     checkNilStr(user.Email),
		FirstName: checkNilStr(user.FirstName),
		LastName:  checkNilStr(user.LastName),
		BirthDate: checkNilTime(user.BirthDate),
	}

	c.JSON(http.StatusOK, userData)
}

func checkNilStr(ptr *string) string {
	if ptr == nil {
		return ""
	}

	return *ptr
}

func checkNilTime(ptr *time.Time) time.Time {
	if ptr == nil {
		return time.Time{}
	}

	return *ptr
}
