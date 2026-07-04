package service

import (
	"chat-gateway/internal/entity"
	"chat-gateway/internal/repository"
	"context"
)

// MessengerService represents interface for work with messenger business logic
type MessengerService interface {
	// CreateUser allows create user into messenger by template
	CreateUser(ctx context.Context, user *entity.User) error
}

type Service struct {
	MessengerService
}

func NewService(messenger repository.MessengerRepository) *Service {
	return &Service{
		MessengerService: messenger,
	}
}
