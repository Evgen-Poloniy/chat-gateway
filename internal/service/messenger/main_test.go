package messenger_test

import (
	"github.com/Evgen-Poloniy/chat-gateway/internal/errs"
	"github.com/Evgen-Poloniy/chat-gateway/internal/service/messenger"
	mock_repository "github.com/Evgen-Poloniy/chat-gateway/internal/service/messenger/mocks"

	"go.uber.org/mock/gomock"
)

var (
	errRepoQueryError = errs.NewAppError(errs.CodeQueryError, "query error", nil)
)

func ptr[T any](v T) *T {
	return &v
}

func setupMockService(ctrl *gomock.Controller) (*messenger.MessengerService, *mock_repository.MockMessengerRepository, *mock_repository.MockMessageBroker, *mock_repository.MockMessengerCache) {
	mockRepo := mock_repository.NewMockMessengerRepository(ctrl)
	mockBroker := mock_repository.NewMockMessageBroker(ctrl)
	mockCache := mock_repository.NewMockMessengerCache(ctrl)

	svc := messenger.NewMessengerService(mockRepo, mockBroker, mockCache)

	return svc, mockRepo, mockBroker, mockCache
}
