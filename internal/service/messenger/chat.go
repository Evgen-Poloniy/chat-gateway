package messenger

import (
	"context"

	"github.com/Evgen-Poloniy/chat-gateway/internal/dto"
	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
)

// CreateDirectChat accept user IDs and create direct chat.
func (m *MessengerService) CreateDirectChat(ctx context.Context, chat *entity.DirectChat) error {
	return m.messengerRepository.CreateDirectChat(ctx, chat)
}

// CreateGroupChat accept user IDs, chat name, chat owner user_id and create group chat between several users.
func (m *MessengerService) CreateGroupChat(ctx context.Context, chat *entity.GroupChat) error {
	return m.messengerRepository.CreateGroupChat(ctx, chat)
}

// SendMessage sends message into target chat.
func (m *MessengerService) SendMessage(message *dto.SendMessageReq) error {
	return m.messageBroker.SendMessage(message)
}
