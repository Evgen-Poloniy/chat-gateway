package resolver_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Evgen-Poloniy/chat-gateway/internal/errs"
	"github.com/Evgen-Poloniy/chat-gateway/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock_repository "github.com/Evgen-Poloniy/chat-gateway/internal/service/resolver/mocks"
)

func TestResolverService_SubscribeOnEventChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc, mockChannel, _ := setupMockService(ctrl)
	ctx := context.Background()

	mockChannel.EXPECT().SubscribeOnEventChannel(ctx).Times(1)

	svc.SubscribeOnEventChannel(ctx)
}

func TestResolverService_ResolveEvent(t *testing.T) {
	testErr := errors.New("test internal error")
	validUUID := uuid.New().String()

	tests := []struct {
		name              string
		mock              func(mockChannel *mock_repository.MockChannel, mockCache *mock_repository.MockCache)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name: "Success",
			mock: func(mockChannel *mock_repository.MockChannel, mockCache *mock_repository.MockCache) {
				mockChannel.EXPECT().
					ResolveEvent(gomock.Any()).
					Return(&model.EventMessage{
						MessageID: "msg-123",
						ChatID:    "chat-123",
						SenderID:  "sender-123",
						Text:      "hello",
						CreatedAt: time.Now(),
					}, nil)

				mockCache.EXPECT().
					GetChatMembers(gomock.Any(), "chat-123").
					Return([]string{validUUID}, nil)

				mockCache.EXPECT().
					ExpireChatID(gomock.Any(), "chat-123").
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Error - Channel.ResolveEvent Failed",
			mock: func(mockChannel *mock_repository.MockChannel, mockCache *mock_repository.MockCache) {
				mockChannel.EXPECT().
					ResolveEvent(gomock.Any()).
					Return(nil, testErr)
			},
			wantErr: true,
		},
		{
			name: "Error - Cache.GetChatMembers Failed",
			mock: func(mockChannel *mock_repository.MockChannel, mockCache *mock_repository.MockCache) {
				mockChannel.EXPECT().
					ResolveEvent(gomock.Any()).
					Return(&model.EventMessage{ChatID: "chat-123"}, nil)

				mockCache.EXPECT().
					GetChatMembers(gomock.Any(), "chat-123").
					Return(nil, testErr)
			},
			wantErr: true,
		},
		{
			name: "Error - Cache.ExpireChatID Failed",
			mock: func(mockChannel *mock_repository.MockChannel, mockCache *mock_repository.MockCache) {
				mockChannel.EXPECT().
					ResolveEvent(gomock.Any()).
					Return(&model.EventMessage{ChatID: "chat-123"}, nil)

				mockCache.EXPECT().
					GetChatMembers(gomock.Any(), "chat-123").
					Return([]string{validUUID}, nil)

				mockCache.EXPECT().
					ExpireChatID(gomock.Any(), "chat-123").
					Return(testErr)
			},
			wantErr: true,
		},
		{
			name: "Error - ParseUUIDs Failed (Empty UserIDs)",
			mock: func(mockChannel *mock_repository.MockChannel, mockCache *mock_repository.MockCache) {
				mockChannel.EXPECT().
					ResolveEvent(gomock.Any()).
					Return(&model.EventMessage{ChatID: "chat-123"}, nil)

				mockCache.EXPECT().
					GetChatMembers(gomock.Any(), "chat-123").
					Return([]string{}, nil)

				mockCache.EXPECT().
					ExpireChatID(gomock.Any(), "chat-123").
					Return(nil)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeEmptyUserIDs,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc, mockChannel, mockCache := setupMockService(ctrl)
			tt.mock(mockChannel, mockCache)

			ctx := context.Background()
			result, err := svc.ResolveEvent(ctx)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)

				var appErr *errs.AppError
				if errors.As(err, &appErr) && tt.expectedErrorCode != "" {
					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "msg-123", result.MessageID)
				assert.Equal(t, "chat-123", result.ChatID)
				assert.Len(t, result.UserIDs, 1)
			}
		})
	}
}
