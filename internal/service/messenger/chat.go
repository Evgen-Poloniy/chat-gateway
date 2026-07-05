package messenger

import (
	"chat-gateway/internal/entity"
	"context"
)

// CreateChat accept user IDs and create direct chat.
func (m *MessengerService) CreateDirectChat(ctx context.Context, chat *entity.DirectChat) error {
	return m.messengerRepository.CreateDirectChat(ctx, chat)
}
