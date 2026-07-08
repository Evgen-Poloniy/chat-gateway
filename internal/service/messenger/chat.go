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

// GetChatsByUserID gets chat by user_id with limits and offset
func (m *MessengerService) GetChatsByUserID(ctx context.Context, userID int64, limit, offset int) ([]entity.Chat, error) {
	return m.messengerRepository.GetChatsByUserID(ctx, userID, limit, offset)
}

// SendMessage sends message into target chat.
func (m *MessengerService) SendMessage(message *dto.SendMessageReq) error {
	return m.messageBroker.SendMessage(message)
}
