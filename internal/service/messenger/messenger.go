package messenger

import (
	"context"

	"github.com/Evgen-Poloniy/chat-gateway/internal/dto"
	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
)

// MessengerRepository represents interface for work with the messenger database.
type MessengerRepository interface {
	// CreateUser allows create user into messenger database.
	CreateUser(ctx context.Context, user *entity.User) error

	// GetUserDataByUsername gets all data about user from the messenger database.
	GetUserDataByUsername(ctx context.Context, username string) (*entity.User, error)

	// UpdateUser updates data about user into messenger database.
	UpdateUser(ctx context.Context, user *entity.UpdateUser) error

	// CreateChat accept user IDs and create direct chat.
	CreateDirectChatReq(ctx context.Context, chat *entity.DirectChat) error

	// CreateChat accept user IDs and create group chat.
	CreateGroupChatReq(ctx context.Context, chat *entity.GroupChat) error
}

// MessageBroker represents interface for work with the messenger broker
type MessageBroker interface {
	// SendMessage sends message into target chat.
	SendMessage(message *dto.SendMessageReq) error
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
