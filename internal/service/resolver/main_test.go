package resolver_test

import (
	"go.uber.org/mock/gomock"

	"github.com/Evgen-Poloniy/chat-gateway/internal/service/resolver"
	mock_repository "github.com/Evgen-Poloniy/chat-gateway/internal/service/resolver/mocks"
)

func setupMockService(ctrl *gomock.Controller) (*resolver.ResolverService, *mock_repository.MockChannel, *mock_repository.MockCache) {
	mockChannel := mock_repository.NewMockChannel(ctrl)
	mockCache := mock_repository.NewMockCache(ctrl)
	svc := resolver.NewResolverService(mockChannel, mockCache)
	return svc, mockChannel, mockCache
}
