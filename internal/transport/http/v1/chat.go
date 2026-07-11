package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Evgen-Poloniy/chat-gateway/internal/dto"
	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
	errs "github.com/Evgen-Poloniy/chat-gateway/pkg/errors"

	"github.com/gin-gonic/gin"
)

// CreateDirectChat create direct chat into the messenger database.
func (h *Handler) CreateDirectChat(c *gin.Context) {
	var req dto.CreateDirectChatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&errs.HttpError{
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
			Message:    err.Error(),
		})
		return
	}

	chat := &entity.DirectChat{
		ParticipantIDs: req.ParticipantIDs,
	}

	if err := h.messenger.CreateDirectChat(c.Request.Context(), chat); err != nil {
		c.Error(err)
		return
	}

	resp := dto.DirectChatResp{
		ChatID:    chat.ChatID,
		CreatedAt: chat.CreatedAt,
		ChatType:  "direct",
	}
	c.JSON(http.StatusCreated, dto.DataResp{Data: resp})
}

// CreateGroupChat create group chat into the messenger database.
func (h *Handler) CreateGroupChat(c *gin.Context) {
	var req dto.CreateGroupChatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&errs.HttpError{
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
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
		c.Error(err)
	}

	resp := dto.GroupChatResp{
		ChatID:         chat.ChatID,
		ParticipantIDs: chat.ParticipantIDs,
		ChatType:       "group",
		Name:           chat.Name,
		Title:          chat.Title,
		Description:    chat.Description,
		CreatedAt:      chat.CreatedAt,
		OwnerID:        chat.OwnerID,
	}
	c.JSON(http.StatusCreated, dto.DataResp{Data: resp})
}

// GetChatsByUserID gets chat by user_id with limits and pages
func (h *Handler) GetChatsByUserID(c *gin.Context) {
	userIdParam := c.Param("user_id")
	if userIdParam == "" {
		c.Error(&errs.HttpError{
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
			Message:    errs.ErrUserIdIsRequired.Error(),
			Err:        errs.ErrUserIdIsRequired,
		})
		return
	}

	userID, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		c.Error(&errs.HttpError{
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
			Message:    "failed to convert parameter 'user_id' to int64",
			Err:        fmt.Errorf("failed to convert parameter 'user_id' to int64: %v", err),
		})
		return
	}

	pageQuery := c.Query("page")
	var page int
	if pageQuery == "" {
		page = 1
	} else {
		var err error
		page, err = strconv.Atoi(pageQuery)
		if err != nil {
			c.Error(&errs.HttpError{
				StatusCode: http.StatusBadRequest,
				Code:       "bad_request",
				Message:    "failed to convert query parameter 'page' to positive int",
				Err:        fmt.Errorf("failed to convert query parameter 'page' to int: %v", err),
			})
			return
		}

		if page < 1 {
			c.Error(&errs.HttpError{
				StatusCode: http.StatusBadRequest,
				Code:       "bad_request",
				Message:    errs.ErrPageRequiredBeGreater.Error(),
				Err:        errs.ErrPageRequiredBeGreater,
			})
			return
		}
	}

	limitQuery := c.Query("limit")
	var limit int
	if limitQuery == "" {
		limit = 100
	} else {
		var err error
		limit, err = strconv.Atoi(limitQuery)
		if err != nil {
			c.Error(&errs.HttpError{
				StatusCode: http.StatusBadRequest,
				Code:       "bad_request",
				Message:    "failed to convert query parameter 'limit' to positive int",
				Err:        fmt.Errorf("failed to convert query parameter 'limit' to int: %v", err),
			})
			return
		}

		if limit < 1 {
			c.Error(&errs.HttpError{
				StatusCode: http.StatusBadRequest,
				Code:       "bad_request",
				Message:    errs.ErrLimitRequiredBeGreater.Error(),
				Err:        errs.ErrLimitRequiredBeGreater,
			})
			return
		}

		if limit > 100 {
			c.Error(&errs.HttpError{
				StatusCode: http.StatusBadRequest,
				Code:       "bad_request",
				Message:    errs.ErrLimitRequiredBeLess.Error(),
				Err:        errs.ErrLimitRequiredBeLess,
			})
			return
		}
	}

	offset := (page - 1) * limit

	chats, err := h.messenger.GetChatsByUserID(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.Error(err)
		return
	}

	if len(chats) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	resp := dto.ChatsResp{
		Chats: chats,
	}

	c.JSON(http.StatusOK, dto.DataResp{Data: resp})
}

// UpdateGroupChat updates data about chat like name, title, description, owner.
func (h *Handler) UpdateGroupChat(c *gin.Context) {
	chatIdParam := c.Param("chat_id")
	if chatIdParam == "" {
		c.Error(&errs.HttpError{
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
			Message:    errs.ErrChatIdIsRequired.Error(),
			Err:        errs.ErrChatIdIsRequired,
		})
		return
	}

	chatID, err := strconv.ParseInt(chatIdParam, 10, 64)
	if err != nil {
		c.Error(&errs.HttpError{
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
			Message:    "failed to convert parameter 'chat_id' to int64",
			Err:        fmt.Errorf("failed to convert parameter 'chat_id' to int64: %v", err),
		})
		return
	}

	var req dto.UpdateGroupChatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&errs.HttpError{
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
			Message:    err.Error(),
			Err:        err,
		})
		return
	}

	chat := entity.UpdateGroupChat{
		ChatID:      chatID,
		Name:        req.Name,
		Title:       req.Title,
		Description: req.Description,
		OwnerID:     req.OwnerID,
	}

	if err := h.messenger.UpdateGroupChat(c.Request.Context(), &chat); err != nil {
		c.Error(err)
	}

	resp := dto.UpdateGroupChatResp{
		UserIDUpdater: chat.UserIDUpdater,
		ChatID:        chat.ChatID,
		ChatType:      "group",
		Name:          *chat.Name,
		Title:         chat.Title,
		Description:   chat.Description,
		CreatedAt:     chat.CreatedAt,
		OwnerID:       *chat.OwnerID,
	}

	c.JSON(http.StatusOK, dto.DataResp{Data: resp})
}

// SendMessage sends message into target chat.
func (h *Handler) SendMessage(c *gin.Context) {
	chatIdParam := c.Param("chat_id")
	if chatIdParam == "" {
		c.Error(&errs.HttpError{
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
			Message:    errs.ErrChatIdIsRequired.Error(),
			Err:        errs.ErrChatIdIsRequired,
		})
		return
	}

	chatID, err := strconv.ParseInt(chatIdParam, 10, 64)
	if err != nil {
		c.Error(&errs.HttpError{
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
			Message:    "failed to convert parameter 'chat_id' to int64",
			Err:        fmt.Errorf("failed to convert parameter 'chat_id' to int64: %v", err),
		})
		return
	}

	var req dto.SendMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&errs.HttpError{
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
			Message:    err.Error(),
		})
		return
	}

	message := entity.SendMessage{
		ChatID:   chatID,
		SenderID: req.SenderID,
		Message:  req.Message,
	}

	if err := h.messenger.SendMessage(c.Request.Context(), &message); err != nil {
		c.Error(err)
		return
	}

	resp := dto.SendMessageResp{
		ChatID:   chatID,
		SenderID: req.SenderID,
		Status:   "sended",
	}

	c.JSON(http.StatusAccepted, dto.DataResp{Data: resp})
}
