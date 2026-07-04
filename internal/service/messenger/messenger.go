package messenger

import "chat-gateway/internal/repository"

// MessengerService represents implementation of messenger interface
type MessengerService struct {
	messenger repository.MessengerRepository
}

func NewMessengerService(messenger repository.MessengerRepository) *MessengerService {
	return &MessengerService{
		messenger: messenger,
	}
}
