package messenger

import (
	"chat-gateway/internal/entity"
	"context"
)

// CreateUser allows create user into messenger by template
func (m *MessengerService) CreateUser(ctx context.Context, user *entity.User) error {
	return m.messenger.CreateUser(ctx, user)
}
