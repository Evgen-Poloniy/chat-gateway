package service

import (
	"chat-gateway/internal/entity"
	"chat-gateway/internal/repository"
	"context"
)

// MessengerService represents interface for work with messenger business logic.
type MessengerService interface {
	// CreateUser allows create user into messenger by template.
	CreateUser(ctx context.Context, user *entity.User) error

	// CreateChat accept user IDs and create direct chat.
	CreateDirectChat(ctx context.Context, chat *entity.DirectChat) error
}

type Service struct {
	MessengerService
}

func NewService(messenger repository.MessengerRepository) *Service {
	return &Service{
		MessengerService: messenger,
	}
}
