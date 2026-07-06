package v1

import (
	"net/http"

	"github.com/Evgen-Poloniy/chat-gateway/internal/dto"
	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
	errs "github.com/Evgen-Poloniy/chat-gateway/pkg/errors"

	"github.com/gin-gonic/gin"
)

// CreateDirectChatReq create direct chat into the messenger database.
func (h *Handler) CreateDirectChatReq(c *gin.Context) {
	var req dto.CreateDirectChatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&errs.HttpError{
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

	if err := h.messenger.CreateDirectChatReq(c.Request.Context(), chat); err != nil {
		c.Error(err)
	}

	resp := &dto.DirectChatResp{
		ChatID:    chat.ChatID,
		CreatedAt: chat.CreatedAt,
	}
	c.JSON(http.StatusCreated, resp)
}

// CreateGroupChatReq create group chat into the messenger database.
func (h *Handler) CreateGroupChatReq(c *gin.Context) {
	var req dto.CreateGroupChatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&errs.HttpError{
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

	if err := h.messenger.CreateGroupChatReq(c.Request.Context(), chat); err != nil {
		c.Error(err)
	}

	resp := &dto.GroupChatResp{
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

// SendMessage sends message into target chat.
func (h *Handler) SendMessage(c *gin.Context) {
	var req dto.SendMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&errs.HttpError{
			StatusCode: http.StatusBadRequest,
			Code:       "BAD_REQUEST",
			Message:    err.Error(),
		})
		return
	}

	if err := h.messenger.SendMessage(&req); err != nil {
		c.Error(err)
		return
	}

	resp := &dto.SendMessageResp{
		ChatID:   req.ChatID,
		SenderID: req.SenderID,
		Status:   "sended",
	}

	c.JSON(http.StatusAccepted, resp)
}
