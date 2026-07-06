package messenger

import (
	"context"

	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
)

// CreateUser allows create user into messenger by template.
func (m *MessengerService) CreateUser(ctx context.Context, user *entity.User) error {
	return m.messengerRepository.CreateUser(ctx, user)
}
