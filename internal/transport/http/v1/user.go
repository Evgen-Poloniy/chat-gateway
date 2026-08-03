package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Evgen-Poloniy/chat-gateway/internal/dto"
	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
	"github.com/Evgen-Poloniy/chat-gateway/internal/errs"
	"github.com/Evgen-Poloniy/chat-gateway/internal/utils"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

// SyncUser processes signup and update-user webhook actions from Casdoor.
func (h *Handler) SingUpUser(c *gin.Context) {
	var webhook dto.CasdoorWebhookReq
	if err := c.ShouldBindJSON(&webhook); err != nil {
		c.Error(errs.NewAppError(errs.CodeBadRequest, err.Error(), err))
		return
	}

	if webhook.Action != "signup" {
		c.Error(errs.NewAppError(errs.CodeBadRequest, "action must be signup", nil))
		return
	}

	var req dto.SignUpUser
	if err := json.Unmarshal(webhook.ExtendedUser, &req); err != nil {
		c.Error(errs.NewAppError(errs.CodeBadRequest, err.Error(), err))
		return
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.Error(errs.NewAppError(errs.CodeBadRequest, "invalid user uuid", err))
		return
	}

	birthDate, err := utils.ParseDate(req.Birthday)
	if err != nil {
		c.Error(errs.NewAppError(errs.CodeBadRequest, "invalid birthday format", err))
		return
	}

	user := &entity.User{
		UserID:    userUUID,
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: birthDate,
		Gender:    req.Gender,
	}

	if err := h.messenger.CreateUser(c.Request.Context(), user); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateUser updates user's information into the messenger database.
func (h *Handler) UpdateUser(c *gin.Context) {
	var webhook dto.CasdoorWebhookReq
	if err := c.ShouldBindJSON(&webhook); err != nil {
		c.Error(errs.NewAppError(errs.CodeBadRequest, err.Error(), err))
		return
	}

	if webhook.Action != "update-user" {
		c.Error(errs.NewAppError(errs.CodeBadRequest, "action must be update-user", nil))
		return
	}

	var req dto.UpdateUser
	if err := json.Unmarshal(webhook.ExtendedUser, &req); err != nil {
		c.Error(errs.NewAppError(errs.CodeBadRequest, err.Error(), err))
		return
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.Error(errs.NewAppError(errs.CodeBadRequest, "invalid user uuid", err))
		return
	}

	birthDate, err := utils.ParseDate(req.Birthday)
	if err != nil {
		c.Error(errs.NewAppError(errs.CodeBadRequest, "invalid birthday format", err))
		return
	}

	user := &entity.UpdateUser{
		UserID:    userUUID,
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: birthDate,
		Gender:    req.Gender,
	}

	if err := h.messenger.UpdateUser(c.Request.Context(), user); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteUser processes delete-user webhook action from Casdoor.
func (h *Handler) DeleteUser(c *gin.Context) {
	var webhook dto.CasdoorWebhookReq
	if err := c.ShouldBindJSON(&webhook); err != nil {
		c.Error(errs.NewAppError(errs.CodeBadRequest, err.Error(), err))
		return
	}

	if webhook.Action != "delete-user" {
		c.Error(errs.NewAppError(errs.CodeBadRequest, "action must be delete-user", nil))
		return
	}

	var req dto.DeleteUser
	if err := json.Unmarshal(webhook.ExtendedUser, &req); err != nil {
		c.Error(errs.NewAppError(errs.CodeBadRequest, err.Error(), err))
		return
	}

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.Error(errs.NewAppError(errs.CodeBadRequest, "invalid user uuid", err))
		return
	}

	if err := h.messenger.DeleteUser(c.Request.Context(), userUUID); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// SearchUser represents searching all data about user from the messenger database.
func (h *Handler) SearchUser(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.Error(errs.NewAppError(
			errs.CodeUsernameIsRequired,
			errs.ErrUsernameIsRequired.Error(),
			errs.ErrUsernameIsRequired,
		))
		return
	}

	user, err := h.messenger.GetUserDataByUsername(c.Request.Context(), username)
	if err != nil {
		c.Error(err)
		return
	}

	resp := dto.UserDataResp{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthDate: user.BirthDate,
		CreatedAt: user.CreatedAt,
		Gender:    user.Gender,
	}

	c.JSON(http.StatusOK, dto.DataResp{Data: resp})
}

// GetUserDataByUserID gets all data about user from the messenger database by user_id.
func (h *Handler) GetUserDataByUserID(c *gin.Context) {
	userIdParam := c.Param("user_id")
	if userIdParam == "" {
		c.Error(errs.NewAppError(
			errs.CodeUserIdIsRequired,
			errs.ErrUserIdIsRequired.Error(),
			errs.ErrUserIdIsRequired,
		))
		return
	}

	userID, err := uuid.Parse(userIdParam)
	if err != nil {
		c.Error(errs.NewAppError(
			errs.CodeBadRequest,
			"failed to parse parameter 'user_id' as uuid",
			fmt.Errorf("failed to parse parameter 'user_id' as uuid: %v", err),
		))
		return
	}

	user, err := h.messenger.GetUserDataByUserID(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	resp := dto.UserDataResp{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthDate: user.BirthDate,
		CreatedAt: user.CreatedAt,
		Gender:    user.Gender,
	}

	c.JSON(http.StatusOK, dto.DataResp{Data: resp})
}
