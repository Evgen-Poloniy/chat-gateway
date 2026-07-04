package v1

import (
	"chat-gateway/internal/dto"
	"chat-gateway/internal/entity"
	errs "chat-gateway/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterUser register user by username and details about user.
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
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: req.BirthDate,
		Gender:    req.Gender,
	}

	if err := h.messenger.CreateUser(c.Request.Context(), user); err != nil {
		c.Error(&errs.AppError{
			StatusCode: http.StatusInternalServerError,
			Code:       "INTERNAL_SERVER_ERROR",
			Message:    "error when create user",
		})
	}

	userData := &dto.UserData{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthDate: user.BirthDate,
		Gender:    user.Gender,
	}

	c.JSON(http.StatusCreated, userData)
}
