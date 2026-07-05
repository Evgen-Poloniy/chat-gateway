package messenger

import (
	"chat-gateway/internal/entity"
	"context"
)

// MessengerRepository represents interface for work with the messenger database.
type MessengerRepository interface {
	// GetUserIdByUsername allows get username from messenger database and get unique used_id.
	GetUserIdByUsername(ctx context.Context, username string) (int64, error)

	// CreateUser allows create user into messenger database.
	CreateUser(ctx context.Context, user *entity.User) error

	// CreateChat accept user IDs and create direct chat.
	CreateDirectChat(ctx context.Context, chat *entity.DirectChat) error
}

// MessengerService represents implementation of messenger interface.
type MessengerService struct {
	messenger MessengerRepository
}

func NewMessengerService(messenger MessengerRepository) *MessengerService {
	return &MessengerService{
		messenger: messenger,
	}
}
