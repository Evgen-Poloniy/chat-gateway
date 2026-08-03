package messenger

import (
	"context"
	"fmt"

	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
	"github.com/Evgen-Poloniy/chat-gateway/internal/errs"
	"github.com/Evgen-Poloniy/chat-gateway/internal/model"
	"github.com/google/uuid"
)

// CreateUser allows create user into messenger by template.
func (m *MessengerService) CreateUser(ctx context.Context, user *entity.User) error {
	if err := m.validateData(ctx, user); err != nil {
		return err
	}

	userModel := model.User{
		UserID:    uuid.New(),
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthDate: user.BirthDate,
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
		Gender:    user.Gender,
	}

	if err := m.messengerRepository.UpdateUser(ctx, &userModel); err != nil {
		return err
	}

	user.Username = userModel.Username
	user.Email = userModel.Email
	user.FirstName = userModel.FirstName
	user.LastName = userModel.LastName
	user.BirthDate = userModel.BirthDate
	user.Gender = userModel.Gender
	user.CreatedAt = userModel.CreatedAt

	return nil
}

// DeleteUser deletes user from the messenger database.
func (s *MessengerService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.messengerRepository.DeleteUser(ctx, userID)
}

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
func (m *MessengerService) GetUserDataByUserID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
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

// GetUserIDsByChatID gets user_id by all users who are in the chat.
func (m *MessengerService) GetUserIDsByChatID(ctx context.Context, chatID uuid.UUID) (uuid.UUIDs, error) {
	cacheUserIDs, err := m.messengerCache.GetChatMembers(ctx, chatID.String())
	if err != nil {
		return nil, err
	}
	if len(cacheUserIDs) == 0 {
		dbUserIDs, err := m.messengerRepository.GetUserIDsByChatID(ctx, chatID)
		if err != nil {
			return nil, err
		}
		if len(dbUserIDs) == 0 {
			return nil, errs.NewAppError(
				errs.CodeChatNotFound,
				fmt.Sprintf("chat %s has no registered members", chatID),
				fmt.Errorf("chat %s has no registered members", chatID),
			)
		}

		return dbUserIDs, nil
	}

	userIDs := make(uuid.UUIDs, 0, len(cacheUserIDs))

	for _, userID := range cacheUserIDs {
		id, err := uuid.Parse(userID)
		if err != nil {
			return nil, errs.NewAppError(
				errs.CodeValidationError,
				"validation error: invalid user uuid",
				fmt.Errorf("validation error: invalid user uuid: %w", err),
			)
		}
		userIDs = append(userIDs, id)
	}

	return userIDs, nil
}
