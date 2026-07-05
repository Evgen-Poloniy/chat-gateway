package v1

import (
	"chat-gateway/internal/entity"
	"context"
)

// MessengerService represents interface for work with messenger business logic.
type MessengerService interface {
	// CreateUser allows create user into messenger by template.
	CreateUser(ctx context.Context, user *entity.User) error

	// CreateChat accept user IDs and create direct chat.
	CreateDirectChat(ctx context.Context, chat *entity.DirectChat) error
}

type Handler struct {
	messenger MessengerService
}

func NewHandler(messenger MessengerService) *Handler {
	return &Handler{
		messenger: messenger,
	}
}
