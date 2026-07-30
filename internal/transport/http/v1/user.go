package v1

import (
	"fmt"
	"net/http"

	"github.com/Evgen-Poloniy/chat-gateway/internal/dto"
	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
	"github.com/Evgen-Poloniy/chat-gateway/internal/errs"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

// RegisterUser register user by username and details about user.
func (h *Handler) RegisterUser(c *gin.Context) {
	var req dto.RegisterUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.NewAppError(
			errs.CodeBadRequest,
			err.Error(),
			err,
		))
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

	c.JSON(http.StatusCreated, dto.DataResp{Data: resp})
}

// // UpdateChat updates data about chat like name, title, description, owner.
// func (h *Handler) UpdateUser(c *gin.Context) {
// 	var req dto.UpdateUserReq
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.Error(errs.NewAppError(
// 			errs.CodeBadRequest,
// 			err.Error(),
// 			err,
// 		))
// 		return
// 	}

// 	user := &entity.UpdateUser{
// 		UserID:    userID,
// 		Username:  req.Username,
// 		Email:     req.Email,
// 		FirstName: req.FirstName,
// 		LastName:  req.LastName,
// 		BirthDate: req.BirthDate,
// 		Gender:    req.Gender,
// 	}

// 	if err := h.messenger.UpdateUser(c.Request.Context(), user); err != nil {
// 		c.Error(err)
// 	}

// 	resp := dto.UserDataResp{
// 		UserID:    user.UserID,
// 		Username:  *user.Username,
// 		Email:     user.Email,
// 		FirstName: user.FirstName,
// 		LastName:  user.LastName,
// 		BirthDate: user.BirthDate,
// 		CreatedAt: user.CreatedAt,
// 		Gender:    user.Gender,
// 	}

// 	c.JSON(http.StatusOK, dto.DataResp{Data: resp})
// }

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
