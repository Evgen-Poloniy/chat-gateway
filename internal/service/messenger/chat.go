package messenger

import (
	"context"

	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
	"github.com/Evgen-Poloniy/chat-gateway/internal/model"
	"github.com/google/uuid"
)

// CreateDirectChat accept user IDs and create direct chat.
func (m *MessengerService) CreateDirectChat(ctx context.Context, chat *entity.DirectChat) error {
	if err := m.validateData(ctx, chat); err != nil {
		return err
	}

	chatModel := model.DirectChat{
		ChatID:         uuid.New(),
		ParticipantIDs: chat.ParticipantIDs,
	}

	if err := m.messengerRepository.CreateDirectChat(ctx, &chatModel); err != nil {
		return err
	}

	chat.ChatID = chatModel.ChatID
	chat.CreatedAt = chatModel.CreatedAt

	if err := m.messengerCache.AddChatMembers(ctx, chat.ChatID.String(), chat.ParticipantIDs.Strings()); err != nil {
		return err
	}

	return nil
}

// CreateGroupChat accept user IDs, chat name, chat owner user_id and create group chat between several users.
func (m *MessengerService) CreateGroupChat(ctx context.Context, chat *entity.GroupChat) error {
	if err := m.validateData(ctx, chat); err != nil {
		return err
	}

	chatModel := model.GroupChat{
		ChatID:      uuid.New(),
		Name:        chat.Name,
		Title:       chat.Title,
		Description: chat.Description,
		OwnerID:     chat.OwnerID,
	}

	if err := m.messengerRepository.CreateGroupChat(ctx, &chatModel); err != nil {
		return err
	}

	chat.ChatID = chatModel.ChatID
	chat.CreatedAt = chatModel.CreatedAt

	if err := m.messengerCache.AddChatMembers(ctx, chat.ChatID.String(), chat.ParticipantIDs.Strings()); err != nil {
		return err
	}

	return nil
}

// GetChatsByUserID gets chat by user_id with limits and offset
func (m *MessengerService) GetChatsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.Chat, error) {
	chatModels, err := m.messengerRepository.GetChatsByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	chats := make([]entity.Chat, 0, len(chatModels))

	for _, chat := range chatModels {
		chats = append(chats, entity.Chat{
			ChatID:      chat.ChatID,
			ChatType:    chat.ChatType,
			Name:        chat.Name,
			Title:       chat.Title,
			Description: chat.Description,
			CreatedAt:   chat.CreatedAt,
			OwnerID:     chat.OwnerID,
		})
	}

	return chats, nil
}

// UpdateGroupChat updates data about chat like name, title, description, owner
func (m *MessengerService) UpdateGroupChat(ctx context.Context, chat *entity.UpdateGroupChat) error {
	if err := m.validateData(ctx, chat); err != nil {
		return err
	}

	chatModel := model.UpdateGroupChat{
		ChatID:      chat.ChatID,
		Name:        chat.Name,
		Title:       chat.Title,
		Description: chat.Description,
		OwnerID:     chat.OwnerID,
	}

	if err := m.messengerRepository.UpdateGroupChat(ctx, &chatModel); err != nil {
		return err
	}

	chat.Name = chatModel.Name
	chat.Title = chatModel.Title
	chat.Description = chatModel.Description
	chat.CreatedAt = chatModel.CreatedAt
	chat.OwnerID = chatModel.OwnerID

	if err := m.messengerCache.ExpireChatID(ctx, chat.ChatID.String()); err != nil {
		return err
	}

	return nil
}

// SendMessage sends message into target chat.
func (m *MessengerService) SendMessage(ctx context.Context, message *entity.SendMessage) error {
	if err := m.validateData(ctx, message); err != nil {
		return err
	}

	messageModel := model.SendMessage{
		ChatID:   message.ChatID,
		SenderID: message.SenderID,
		Message:  message.Message,
	}

	return m.messageBroker.SendMessage(ctx, &messageModel)
}
