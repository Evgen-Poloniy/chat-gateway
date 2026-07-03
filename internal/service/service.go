package service

import "chat-gateway/internal/repository"

// MessengerService represents interface for work with messenger business logic
type MessengerService interface {
}

type Service struct {
	MessengerService
}

func NewService(repository repository.MessengerRepository) *Service {
	return &Service{
		MessengerService: repository,
	}
}
