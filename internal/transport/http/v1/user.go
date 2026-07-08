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

// RegisterUser register user by username and details about user.
func (h *Handler) RegisterUser(c *gin.Context) {
	var req dto.RegisterUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&errs.HttpError{
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
			Message:    err.Error(),
			Err:        err,
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
		c.Error(err)
	}

	resp := &dto.UserDataResp{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthDate: user.BirthDate,
		Gender:    user.Gender,
	}

	c.JSON(http.StatusCreated, resp)
}

// RegisterUser register user by username and details about user.
func (h *Handler) UpdateUser(c *gin.Context) {
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

	var req dto.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&errs.HttpError{
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
			Message:    err.Error(),
			Err:        err,
		})
		return
	}

	user := &entity.UpdateUser{
		UserID:    userID,
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: req.BirthDate,
		Gender:    req.Gender,
	}

	if err := h.messenger.UpdateUser(c.Request.Context(), user); err != nil {
		c.Error(err)
	}

	resp := &dto.UserDataResp{
		UserID:    user.UserID,
		Username:  *user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthDate: user.BirthDate,
		Gender:    user.Gender,
	}

	c.JSON(http.StatusOK, resp)
}

// GetUserDataByUsername represents searching all data about user from the messenger database.
func (h *Handler) GetUserDataByUsername(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.Error(&errs.HttpError{
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
			Message:    errs.ErrUsernameIsRequired.Error(),
			Err:        errs.ErrUsernameIsRequired,
		})
		return
	}

	user, err := h.messenger.GetUserDataByUsername(c.Request.Context(), username)
	if err != nil {
		c.Error(err)
		return
	}

	resp := &dto.UserDataResp{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthDate: user.BirthDate,
		Gender:    user.Gender,
	}

	c.JSON(http.StatusOK, resp)
}

// GetUserDataByUserID gets all data about user from the messenger database by user_id.
func (h *Handler) GetUserDataByUserID(c *gin.Context) {
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

	user, err := h.messenger.GetUserDataByUserID(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
	}

	resp := &dto.UserDataResp{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthDate: user.BirthDate,
		Gender:    user.Gender,
	}

	c.JSON(http.StatusOK, resp)
}
