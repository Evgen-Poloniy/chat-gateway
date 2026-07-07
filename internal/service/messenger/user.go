package messenger

import (
	"context"

	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
)

// GetUserDataByUsername gets all data about user from the messenger database.
func (m *MessengerService) GetUserDataByUsername(ctx context.Context, username string) (*entity.User, error) {
	return m.messengerRepository.GetUserDataByUsername(ctx, username)
}

// CreateUser allows create user into messenger by template.
func (m *MessengerService) CreateUser(ctx context.Context, user *entity.User) error {
	return m.messengerRepository.CreateUser(ctx, user)
}

// UpdateUser updates data about user into messenger database.
func (m *MessengerService) UpdateUser(ctx context.Context, user *entity.UpdateUser) error {
	return m.messengerRepository.UpdateUser(ctx, user)
}
