package messenger

import "chat-gateway/internal/repository"

// MessengerService represents implementation of messenger interface
type MessengerService struct {
	repository repository.MessengerRepository
}

func NewMessengerService(repository repository.MessengerRepository) *MessengerService {
	return &MessengerService{
		repository: repository,
	}
}
