package messenger

import (
	"context"

	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
)

// MessengerRepository represents interface for work with the messenger database.
type MessengerRepository interface {
	// GetUserIdByUsername allows get username from messenger database and get unique used_id.
	GetUserIdByUsername(ctx context.Context, username string) (int64, error)

	// CreateUser allows create user into messenger database.
	CreateUser(ctx context.Context, user *entity.User) error

	// CreateChat accept user IDs and create direct chat.
	CreateDirectChat(ctx context.Context, chat *entity.DirectChat) error

	// CreateChat accept user IDs and create group chat.
	CreateGroupChat(ctx context.Context, chat *entity.GroupChat) error
}

// MessageBroker represents interface for work with the messenger broker
type MessageBroker interface {
}

// MessengerService represents implementation of messenger interface.
type MessengerService struct {
	messengerRepository MessengerRepository
	messageBroker       MessageBroker
}

func NewMessengerService(messengerRepository MessengerRepository, messageBroker MessageBroker) *MessengerService {
	return &MessengerService{
		messengerRepository: messengerRepository,
		messageBroker:       messageBroker,
	}
}
