package v1

import (
	"chat-gateway/internal/entity"
	"context"
)

// MessengerService represents interface for work with messenger business logic.
type MessengerService interface {
	// CreateUser allows create user into messenger by template.
	CreateUser(ctx context.Context, user *entity.User) error

	// CreateDirectChat accept user IDs and create direct chat between two users.
	CreateDirectChat(ctx context.Context, chat *entity.DirectChat) error

	// CreateGroupChat accept user IDs, chat name, chat owner user_id and create group chat between several users.
	CreateGroupChat(ctx context.Context, chat *entity.GroupChat) error
}

type Handler struct {
	messenger MessengerService
}

func NewHandler(messenger MessengerService) *Handler {
	return &Handler{
		messenger: messenger,
	}
}
