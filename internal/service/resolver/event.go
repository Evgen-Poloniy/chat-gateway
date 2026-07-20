package resolver

import (
	"context"

	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
)

// SubscribeToEvents subscribes on broker channel once fro all time of work application.
func (s *ResolverService) SubscribeToEvents(ctx context.Context) error {
	return s.channel.SubscribeToEventChannel(ctx)
}

// ResolveEvent returns event from channel.
func (s *ResolverService) ResolveEvent(ctx context.Context) (*entity.EventMessage, error) {
	event, err := s.channel.ResolveEvent(ctx)
	if err != nil {
		return nil, err
	}

	message := entity.EventMessage{
		MessageID: event.MessageID,
		ChatID:    event.ChatID,
		SenderID:  event.SenderID,
		Text:      event.Text,
		CreatedAt: event.CreatedAt,
	}

	userIDs, err := s.cache.GetChatMembers(ctx, event.ChatID)
	if err != nil {
		return nil, err
	}

	// exception sender user_id from user ids without saving order
	for i, userID := range userIDs {
		if userID == event.SenderID {
			userIDs[i] = userIDs[len(userIDs)-1]
			userIDs = userIDs[:len(userIDs)-1]
			break
		}
	}

	if err := s.cache.ExpireChatID(ctx, event.ChatID); err != nil {
		return nil, err
	}

	if err := message.ParseUUIDs(userIDs); err != nil {
		return nil, err
	}

	return &message, nil
}
