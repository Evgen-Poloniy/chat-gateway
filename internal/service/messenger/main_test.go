package messenger_test

import (
	"github.com/Evgen-Poloniy/chat-gateway/internal/service/messenger"
	mock_repository "github.com/Evgen-Poloniy/chat-gateway/internal/service/messenger/mocks"

	"go.uber.org/mock/gomock"
)

func ptr[T any](v T) *T {
	return &v
}

func setupMockService(ctrl *gomock.Controller) (*messenger.MessengerService, *mock_repository.MockMessengerRepository, *mock_repository.MockMessageBroker) {
	mockRepo := mock_repository.NewMockMessengerRepository(ctrl)
	mockBroker := mock_repository.NewMockMessageBroker(ctrl)

	svc := messenger.NewMessengerService(mockRepo, mockBroker)

	return svc, mockRepo, mockBroker
}
