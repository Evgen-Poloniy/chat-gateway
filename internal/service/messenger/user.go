package messenger

import (
	"context"

	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
	"github.com/Evgen-Poloniy/chat-gateway/internal/model"
)

// GetUserDataByUsername represents searching all data about user from the messenger database.
func (m *MessengerService) GetUserDataByUsername(ctx context.Context, username string) (*entity.User, error) {
	modelUser, err := m.messengerRepository.GetUserDataByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		UserID:    modelUser.UserID,
		Username:  modelUser.Username,
		Email:     modelUser.Email,
		FirstName: modelUser.FirstName,
		LastName:  modelUser.LastName,
		BirthDate: modelUser.BirthDate,
		CreatedAt: modelUser.CreatedAt,
		Gender:    modelUser.Gender,
	}

	return user, nil
}

// GetUserDataByUserID gets all data about user from the messenger database by user_id.
func (m *MessengerService) GetUserDataByUserID(ctx context.Context, userID int64) (*entity.User, error) {
	modelUser, err := m.messengerRepository.GetUserDataByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		UserID:    modelUser.UserID,
		Username:  modelUser.Username,
		Email:     modelUser.Email,
		FirstName: modelUser.FirstName,
		LastName:  modelUser.LastName,
		BirthDate: modelUser.BirthDate,
		CreatedAt: modelUser.CreatedAt,
		Gender:    modelUser.Gender,
	}

	return user, nil
}

// CreateUser allows create user into messenger by template.
func (m *MessengerService) CreateUser(ctx context.Context, user *entity.User) error {
	if err := m.validateData(ctx, user); err != nil {
		return err
	}

	userModel := model.User{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthDate: user.BirthDate,
		CreatedAt: user.CreatedAt,
		Gender:    user.Gender,
	}

	if err := m.messengerRepository.CreateUser(ctx, &userModel); err != nil {
		return err
	}

	user.UserID = userModel.UserID
	user.CreatedAt = userModel.CreatedAt

	return nil
}

// UpdateUser updates data about user into messenger database.
func (m *MessengerService) UpdateUser(ctx context.Context, user *entity.UpdateUser) error {
	if err := m.validateData(ctx, user); err != nil {
		return err
	}

	userModel := model.UpdateUser{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthDate: user.BirthDate,
		CreatedAt: user.CreatedAt,
		Gender:    user.Gender,
	}

	if err := m.messengerRepository.UpdateUser(ctx, &userModel); err != nil {
		return err
	}

	user.CreatedAt = userModel.CreatedAt

	return nil
}
