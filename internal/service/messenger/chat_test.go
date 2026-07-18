package messenger_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
	"github.com/Evgen-Poloniy/chat-gateway/internal/model"
	mock_repository "github.com/Evgen-Poloniy/chat-gateway/internal/service/messenger/mocks"
	errs "github.com/Evgen-Poloniy/chat-gateway/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	errRepoChatNotFound = errs.NewAppError(errs.CodeChatNotFound, "chat not found", nil)
	errBrokerError      = errs.NewAppError(errs.CodeQueryError, "broker error", nil)
)

func TestMessengerService_CreateDirectChat(t *testing.T) {
	tests := []struct {
		name              string
		input             *entity.DirectChat
		mock              func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name: "Success",
			input: &entity.DirectChat{
				ParticipantIDs: []uuid.UUID{uuid.New(), uuid.New()},
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockRepo.EXPECT().
					CreateDirectChat(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, m *model.DirectChat) error {
						m.ChatID = uuid.New()
						m.CreatedAt = time.Now()
						return nil
					})

				mockCache.EXPECT().
					AddChatMembers(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Error - Repository Failed",
			input: &entity.DirectChat{
				ParticipantIDs: []uuid.UUID{uuid.New(), uuid.New()},
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockRepo.EXPECT().
					CreateDirectChat(gomock.Any(), gomock.Any()).
					Return(errRepoQueryError)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - Validation (Missing ParticipantIDs)",
			input: &entity.DirectChat{
				ParticipantIDs: nil,
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Invalid ParticipantIDs count - 1)",
			input: &entity.DirectChat{
				ParticipantIDs: []uuid.UUID{uuid.New()},
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Invalid ParticipantIDs count - 3)",
			input: &entity.DirectChat{
				ParticipantIDs: []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc, mockRepo, _, mockCache := setupMockService(ctrl)

			tt.mock(mockRepo, mockCache)

			err := svc.CreateDirectChat(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				var appErr *errs.AppError
				if errors.As(err, &appErr) {
					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, tt.input.ChatID)
				assert.False(t, tt.input.CreatedAt.IsZero())
			}
		})
	}
}

func TestMessengerService_CreateGroupChat(t *testing.T) {
	now := time.Now()
	generatedChatID := uuid.New()

	tests := []struct {
		name              string
		input             *entity.GroupChat
		mock              func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name: "Success",
			input: &entity.GroupChat{
				ParticipantIDs: []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
				Name:           "Dev Team",
				OwnerID:        uuid.New(),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockRepo.EXPECT().
					CreateGroupChat(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, m *model.GroupChat) error {
						m.ChatID = generatedChatID
						m.CreatedAt = now
						return nil
					})

				mockCache.EXPECT().
					AddChatMembers(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Success - Boundary Values",
			input: &entity.GroupChat{
				ParticipantIDs: []uuid.UUID{uuid.New(), uuid.New()},
				Name:           strings.Repeat("A", 64),
				Title:          ptr(strings.Repeat("B", 64)),
				Description:    ptr(strings.Repeat("C", 255)),
				OwnerID:        uuid.New(),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockRepo.EXPECT().
					CreateGroupChat(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, m *model.GroupChat) error {
						m.ChatID = generatedChatID
						m.CreatedAt = now
						return nil
					})

				mockCache.EXPECT().
					AddChatMembers(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Error - Repository Failed",
			input: &entity.GroupChat{
				ParticipantIDs: []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
				Name:           "Dev Team",
				OwnerID:        uuid.New(),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockRepo.EXPECT().
					CreateGroupChat(gomock.Any(), gomock.Any()).
					Return(errRepoQueryError)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - Validation (Missing Name - empty)",
			input: &entity.GroupChat{
				ParticipantIDs: []uuid.UUID{uuid.New(), uuid.New()},
				Name:           "",
				OwnerID:        uuid.New(),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Name too long - 65)",
			input: &entity.GroupChat{
				ParticipantIDs: []uuid.UUID{uuid.New(), uuid.New()},
				Name:           strings.Repeat("A", 65),
				OwnerID:        uuid.New(),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Title too long - 65)",
			input: &entity.GroupChat{
				ParticipantIDs: []uuid.UUID{uuid.New(), uuid.New()},
				Name:           "Dev Team",
				Title:          ptr(strings.Repeat("A", 65)),
				OwnerID:        uuid.New(),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Description too long - 256)",
			input: &entity.GroupChat{
				ParticipantIDs: []uuid.UUID{uuid.New(), uuid.New()},
				Name:           "Dev Team",
				Description:    ptr(strings.Repeat("A", 256)),
				OwnerID:        uuid.New(),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Missing OwnerID - Nil UUID)",
			input: &entity.GroupChat{
				ParticipantIDs: []uuid.UUID{uuid.New(), uuid.New()},
				Name:           "Dev Team",
				OwnerID:        uuid.Nil,
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Less than 2 participants - 1)",
			input: &entity.GroupChat{
				ParticipantIDs: []uuid.UUID{uuid.New()},
				Name:           "Dev Team",
				OwnerID:        uuid.New(),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Empty participants - 0)",
			input: &entity.GroupChat{
				ParticipantIDs: []uuid.UUID{},
				Name:           "Dev Team",
				OwnerID:        uuid.New(),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc, mockRepo, _, mockCache := setupMockService(ctrl)

			tt.mock(mockRepo, mockCache)

			err := svc.CreateGroupChat(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				var appErr *errs.AppError
				if errors.As(err, &appErr) {
					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, generatedChatID, tt.input.ChatID)
				assert.Equal(t, now, tt.input.CreatedAt)
			}
		})
	}
}

func TestMessengerService_GetChatsByUserID(t *testing.T) {
	userID := uuid.New()
	chatID := uuid.New()
	ownerID := uuid.New()

	tests := []struct {
		name              string
		userID            uuid.UUID
		limit             int
		offset            int
		mock              func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name:   "Success",
			userID: userID,
			limit:  10,
			offset: 0,
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockRepo.EXPECT().
					GetChatsByUserID(gomock.Any(), userID, 10, 0).
					Return([]model.Chat{
						{ChatID: chatID, ChatType: "direct", Name: ptr("Chat 1"), OwnerID: ptr(ownerID)},
					}, nil)
			},
			wantErr: false,
		},
		{
			name:   "Success - Empty Result",
			userID: uuid.New(),
			limit:  10,
			offset: 0,
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockRepo.EXPECT().
					GetChatsByUserID(gomock.Any(), gomock.Any(), 10, 0).
					Return([]model.Chat{}, nil)
			},
			wantErr: false,
		},
		{
			name:   "Error - Repository Failed",
			userID: userID,
			limit:  10,
			offset: 0,
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockRepo.EXPECT().
					GetChatsByUserID(gomock.Any(), userID, 10, 0).
					Return(nil, errRepoQueryError)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc, mockRepo, _, mockCache := setupMockService(ctrl)

			tt.mock(mockRepo, mockCache)

			chats, err := svc.GetChatsByUserID(context.Background(), tt.userID, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				var appErr *errs.AppError
				if errors.As(err, &appErr) {
					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
				}
				assert.Nil(t, chats)
			} else {
				assert.NoError(t, err)
				if tt.name == "Success" {
					require.Len(t, chats, 1)
					assert.Equal(t, chatID, chats[0].ChatID)
				}
			}
		})
	}
}

func TestMessengerService_UpdateGroupChat(t *testing.T) {
	now := time.Now()
	chatID := uuid.New()
	updaterID := uuid.New()
	newOwnerID := uuid.New()

	tests := []struct {
		name              string
		input             *entity.UpdateGroupChat
		mock              func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name: "Success",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: updaterID,
				ChatID:        chatID,
				Name:          ptr("New Name"),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockRepo.EXPECT().
					UpdateGroupChat(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, m *model.UpdateGroupChat) error {
						m.CreatedAt = now
						return nil
					})

				mockCache.EXPECT().
					ExpireChatID(gomock.Any(), chatID.String()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Success - Boundary Values",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: updaterID,
				ChatID:        chatID,
				Name:          ptr(strings.Repeat("A", 64)),
				Title:         ptr(strings.Repeat("B", 64)),
				Description:   ptr(strings.Repeat("C", 255)),
				OwnerID:       ptr(newOwnerID),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockRepo.EXPECT().
					UpdateGroupChat(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, m *model.UpdateGroupChat) error {
						m.CreatedAt = now
						return nil
					})

				mockCache.EXPECT().
					ExpireChatID(gomock.Any(), chatID.String()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Error - Repository Failed",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: updaterID,
				ChatID:        chatID,
				Name:          ptr("New Name"),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockRepo.EXPECT().
					UpdateGroupChat(gomock.Any(), gomock.Any()).
					Return(errRepoQueryError)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - Validation (Missing UserIDUpdater - Nil UUID)",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: uuid.Nil,
				ChatID:        chatID,
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Missing ChatID - Nil UUID)",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: updaterID,
				ChatID:        uuid.Nil,
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Name too long - 65)",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: updaterID,
				ChatID:        chatID,
				Name:          ptr(strings.Repeat("A", 65)),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Name empty string)",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: updaterID,
				ChatID:        chatID,
				Name:          ptr(""),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Title too long - 65)",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: updaterID,
				ChatID:        chatID,
				Title:         ptr(strings.Repeat("A", 65)),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Description too long - 256)",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: updaterID,
				ChatID:        chatID,
				Description:   ptr(strings.Repeat("A", 256)),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc, mockRepo, _, mockCache := setupMockService(ctrl)

			tt.mock(mockRepo, mockCache)

			err := svc.UpdateGroupChat(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				var appErr *errs.AppError
				if errors.As(err, &appErr) {
					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, now, tt.input.CreatedAt)
			}
		})
	}
}

func TestMessengerService_SendMessage(t *testing.T) {
	tests := []struct {
		name              string
		input             *entity.SendMessage
		mock              func(mockBroker *mock_repository.MockMessageBroker)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name: "Success",
			input: &entity.SendMessage{
				ChatID:   uuid.New(),
				SenderID: uuid.New(),
				Message:  "Hello World",
			},
			mock: func(mockBroker *mock_repository.MockMessageBroker) {
				mockBroker.EXPECT().
					SendMessage(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Success - Boundary Values",
			input: &entity.SendMessage{
				ChatID:   uuid.New(),
				SenderID: uuid.New(),
				Message:  strings.Repeat("A", 4096),
			},
			mock: func(mockBroker *mock_repository.MockMessageBroker) {
				mockBroker.EXPECT().
					SendMessage(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Error - Broker Failed",
			input: &entity.SendMessage{
				ChatID:   uuid.New(),
				SenderID: uuid.New(),
				Message:  "Hello World",
			},
			mock: func(mockBroker *mock_repository.MockMessageBroker) {
				mockBroker.EXPECT().
					SendMessage(gomock.Any(), gomock.Any()).
					Return(errBrokerError)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - Validation (Missing ChatID - Nil UUID)",
			input: &entity.SendMessage{
				ChatID:   uuid.Nil,
				SenderID: uuid.New(),
				Message:  "Hello World",
			},
			mock:              func(mockBroker *mock_repository.MockMessageBroker) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Missing SenderID - Nil UUID)",
			input: &entity.SendMessage{
				ChatID:   uuid.New(),
				SenderID: uuid.Nil,
				Message:  "Hello World",
			},
			mock:              func(mockBroker *mock_repository.MockMessageBroker) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Empty Message)",
			input: &entity.SendMessage{
				ChatID:   uuid.New(),
				SenderID: uuid.New(),
				Message:  "",
			},
			mock:              func(mockBroker *mock_repository.MockMessageBroker) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Message too long - 4097)",
			input: &entity.SendMessage{
				ChatID:   uuid.New(),
				SenderID: uuid.New(),
				Message:  strings.Repeat("A", 4097),
			},
			mock:              func(mockBroker *mock_repository.MockMessageBroker) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc, _, mockBroker, _ := setupMockService(ctrl)

			tt.mock(mockBroker)

			err := svc.SendMessage(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				var appErr *errs.AppError
				if errors.As(err, &appErr) {
					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
