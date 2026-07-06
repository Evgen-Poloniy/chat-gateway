package v1

import (
	"fmt"
	"net/http"

	"github.com/Evgen-Poloniy/chat-gateway/internal/dto"
	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
	errs "github.com/Evgen-Poloniy/chat-gateway/pkg/errors"

	"github.com/gin-gonic/gin"
)

// CreateDirectChat create direct chat into the messenger database.
func (h *Handler) CreateDirectChat(c *gin.Context) {
	var req dto.CreateDirectChat
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&errs.AppError{
			StatusCode: http.StatusBadRequest,
			Code:       "BAD_REQUEST",
			Message:    err.Error(),
		})
		return
	}

	chat := &entity.DirectChat{
		SenderID:    req.SenderID,
		RecipientID: req.RecipientID,
	}

	if err := h.messenger.CreateDirectChat(c.Request.Context(), chat); err != nil {
		c.Error(&errs.AppError{
			StatusCode: http.StatusBadRequest,
			Code:       "INTERNAL_SERVER_ERROR",
			Message:    fmt.Sprintf("database error: %v", err),
		})
	}

	resp := &dto.DirectChatData{
		ChatID:    chat.ChatID,
		CreatedAt: chat.CreatedAt,
	}
	c.JSON(http.StatusCreated, resp)
}

// CreateGroupChat create group chat into the messenger database.
func (h *Handler) CreateGroupChat(c *gin.Context) {
	var req dto.CreateGroupChat
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&errs.AppError{
			StatusCode: http.StatusBadRequest,
			Code:       "BAD_REQUEST",
			Message:    err.Error(),
		})
		return
	}

	chat := &entity.GroupChat{
		ParticipantIDs: req.ParticipantIDs,
		Name:           req.Name,
		Title:          req.Title,
		Description:    req.Description,
		OwnerID:        req.OwnerID,
	}

	if err := h.messenger.CreateGroupChat(c.Request.Context(), chat); err != nil {
		c.Error(&errs.AppError{
			StatusCode: http.StatusBadRequest,
			Code:       "INTERNAL_SERVER_ERROR",
			Message:    fmt.Sprintf("database error: %v", err),
		})
	}

	resp := &dto.GroupChatData{
		ChatID:         chat.ChatID,
		ParticipantIDs: chat.ParticipantIDs,
		Name:           chat.Name,
		Title:          chat.Title,
		Description:    chat.Description,
		CreatedAt:      chat.CreatedAt,
		OwnerID:        chat.OwnerID,
	}
	c.JSON(http.StatusCreated, resp)
}
