package messenger

import (
	"context"

	"github.com/Evgen-Poloniy/chat-gateway/internal/dto"
	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
)

// CreateDirectChatReq accept user IDs and create direct chat.
func (m *MessengerService) CreateDirectChatReq(ctx context.Context, chat *entity.DirectChat) error {
	return m.messengerRepository.CreateDirectChatReq(ctx, chat)
}

// CreateGroupChatReq accept user IDs, chat name, chat owner user_id and create group chat between several users.
func (m *MessengerService) CreateGroupChatReq(ctx context.Context, chat *entity.GroupChat) error {
	return m.messengerRepository.CreateGroupChatReq(ctx, chat)
}

// SendMessage sends message into target chat.
func (m *MessengerService) SendMessage(message *dto.SendMessageReq) error {
	return m.messageBroker.SendMessage(message)
}
