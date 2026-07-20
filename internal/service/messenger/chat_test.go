package messenger_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
	"github.com/Evgen-Poloniy/chat-gateway/internal/errs"
	"github.com/Evgen-Poloniy/chat-gateway/internal/model"
	mock_repository "github.com/Evgen-Poloniy/chat-gateway/internal/service/messenger/mocks"
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
	chatID1 := uuid.New()
	chatID2 := uuid.New()
	ownerID := uuid.New()
	now := time.Now()

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
						{
							ChatID:      chatID1,
							ChatType:    "direct",
							Name:        ptr("Chat 1"),
							Title:       nil,
							Description: nil,
							CreatedAt:   now,
							OwnerID:     &ownerID,
						},
						{
							ChatID:      chatID2,
							ChatType:    "group",
							Name:        ptr("Group Chat"),
							Title:       ptr("Title"),
							Description: ptr("Desc"),
							CreatedAt:   now,
							OwnerID:     &userID,
						},
					}, nil)
			},
			wantErr: false,
		},
		{
			name:   "Error - Chat Not Found (Empty Result)",
			userID: userID,
			limit:  10,
			offset: 0,
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockRepo.EXPECT().
					GetChatsByUserID(gomock.Any(), userID, 10, 0).
					Return([]model.Chat{}, nil)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeChatNotFound,
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
				if assert.ErrorAs(t, err, &appErr) {
					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
				}
				assert.Nil(t, chats)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, chats)
				assert.Len(t, chats, 2)

				assert.Equal(t, chatID1, chats[0].ChatID)
				assert.Equal(t, "direct", chats[0].ChatType)
				assert.Equal(t, "Chat 1", *chats[0].Name)

				assert.Equal(t, chatID2, chats[1].ChatID)
				assert.Equal(t, "group", chats[1].ChatType)
				assert.Equal(t, "Group Chat", *chats[1].Name)
			}
		})
	}
}

func TestMessengerService_UpdateGroupChat(t *testing.T) {
	now := time.Now()
	chatID := uuid.New()
	updaterID := uuid.New()

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
				mockCache.EXPECT().
					IsUserInChat(gomock.Any(), chatID.String(), updaterID.String()).
					Return(true, nil)

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
			name: "Success - User not in cache but in database",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: updaterID,
				ChatID:        chatID,
				Name:          ptr("New Name"),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockCache.EXPECT().
					IsUserInChat(gomock.Any(), chatID.String(), updaterID.String()).
					Return(false, nil)

				mockRepo.EXPECT().
					IsUserInChat(gomock.Any(), chatID, updaterID).
					Return(true, nil)

				mockCache.EXPECT().
					AddChatMembers(gomock.Any(), chatID.String(), []string{updaterID.String()}).
					Return(nil)

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
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockCache.EXPECT().
					IsUserInChat(gomock.Any(), chatID.String(), updaterID.String()).
					Return(true, nil)

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
			name: "Success - Only Description Update",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: updaterID,
				ChatID:        chatID,
				Description:   ptr("New description"),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockCache.EXPECT().
					IsUserInChat(gomock.Any(), chatID.String(), updaterID.String()).
					Return(true, nil)

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
			name: "Success - Only Title Update",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: updaterID,
				ChatID:        chatID,
				Title:         ptr("New Title"),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockCache.EXPECT().
					IsUserInChat(gomock.Any(), chatID.String(), updaterID.String()).
					Return(true, nil)

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
				mockCache.EXPECT().
					IsUserInChat(gomock.Any(), chatID.String(), updaterID.String()).
					Return(true, nil)

				mockRepo.EXPECT().
					UpdateGroupChat(gomock.Any(), gomock.Any()).
					Return(errRepoQueryError)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - Cache Check Failed",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: updaterID,
				ChatID:        chatID,
				Name:          ptr("New Name"),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockCache.EXPECT().
					IsUserInChat(gomock.Any(), chatID.String(), updaterID.String()).
					Return(false, errCacheFailed)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeRedisError,
		},
		{
			name: "Error - Database Check Failed",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: updaterID,
				ChatID:        chatID,
				Name:          ptr("New Name"),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockCache.EXPECT().
					IsUserInChat(gomock.Any(), chatID.String(), updaterID.String()).
					Return(false, nil)

				mockRepo.EXPECT().
					IsUserInChat(gomock.Any(), chatID, updaterID).
					Return(false, errRepoQueryError)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - User Not In Chat",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: updaterID,
				ChatID:        chatID,
				Name:          ptr("New Name"),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockCache.EXPECT().
					IsUserInChat(gomock.Any(), chatID.String(), updaterID.String()).
					Return(false, nil)

				mockRepo.EXPECT().
					IsUserInChat(gomock.Any(), chatID, updaterID).
					Return(false, nil)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeUserIsNotInChat,
		},
		{
			name: "Error - Add Chat Members to Cache Failed",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: updaterID,
				ChatID:        chatID,
				Name:          ptr("New Name"),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockCache.EXPECT().
					IsUserInChat(gomock.Any(), chatID.String(), updaterID.String()).
					Return(false, nil)

				mockRepo.EXPECT().
					IsUserInChat(gomock.Any(), chatID, updaterID).
					Return(true, nil)

				mockCache.EXPECT().
					AddChatMembers(gomock.Any(), chatID.String(), []string{updaterID.String()}).
					Return(errCacheFailed)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeRedisError,
		},
		{
			name: "Error - Expire Cache Failed",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: updaterID,
				ChatID:        chatID,
				Name:          ptr("New Name"),
			},
			mock: func(mockRepo *mock_repository.MockMessengerRepository, mockCache *mock_repository.MockMessengerCache) {
				mockCache.EXPECT().
					IsUserInChat(gomock.Any(), chatID.String(), updaterID.String()).
					Return(true, nil)

				mockRepo.EXPECT().
					UpdateGroupChat(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, m *model.UpdateGroupChat) error {
						m.CreatedAt = now
						return nil
					})

				mockCache.EXPECT().
					ExpireChatID(gomock.Any(), chatID.String()).
					Return(errCacheFailed)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeRedisError,
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
	chatID := uuid.New()
	senderID := uuid.New()

	tests := []struct {
		name              string
		input             *entity.SendMessage
		mock              func(mockBroker *mock_repository.MockMessageBroker, mockCache *mock_repository.MockMessengerCache, mockRepo *mock_repository.MockMessengerRepository)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name: "Success - User In Cache",
			input: &entity.SendMessage{
				ChatID:   chatID,
				SenderID: senderID,
				Message:  "Hello World",
			},
			mock: func(mockBroker *mock_repository.MockMessageBroker, mockCache *mock_repository.MockMessengerCache, mockRepo *mock_repository.MockMessengerRepository) {
				mockCache.EXPECT().
					IsUserInChat(gomock.Any(), chatID.String(), senderID.String()).
					Return(true, nil)

				mockBroker.EXPECT().
					SendMessage(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Success - User Not In Cache, Found In DB",
			input: &entity.SendMessage{
				ChatID:   chatID,
				SenderID: senderID,
				Message:  "Hello World",
			},
			mock: func(mockBroker *mock_repository.MockMessageBroker, mockCache *mock_repository.MockMessengerCache, mockRepo *mock_repository.MockMessengerRepository) {
				mockCache.EXPECT().
					IsUserInChat(gomock.Any(), chatID.String(), senderID.String()).
					Return(false, nil)

				mockRepo.EXPECT().
					IsUserInChat(gomock.Any(), chatID, senderID).
					Return(true, nil)

				mockCache.EXPECT().
					AddChatMembers(gomock.Any(), chatID.String(), []string{senderID.String()}).
					Return(nil)

				mockBroker.EXPECT().
					SendMessage(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Error - User Not In Chat (DB)",
			input: &entity.SendMessage{
				ChatID:   chatID,
				SenderID: senderID,
				Message:  "Hello World",
			},
			mock: func(mockBroker *mock_repository.MockMessageBroker, mockCache *mock_repository.MockMessengerCache, mockRepo *mock_repository.MockMessengerRepository) {
				mockCache.EXPECT().
					IsUserInChat(gomock.Any(), chatID.String(), senderID.String()).
					Return(false, nil)

				mockRepo.EXPECT().
					IsUserInChat(gomock.Any(), chatID, senderID).
					Return(false, nil)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeUserIsNotInChat,
		},
		{
			name: "Error - Broker Failed",
			input: &entity.SendMessage{
				ChatID:   chatID,
				SenderID: senderID,
				Message:  "Hello World",
			},
			mock: func(mockBroker *mock_repository.MockMessageBroker, mockCache *mock_repository.MockMessengerCache, mockRepo *mock_repository.MockMessengerRepository) {
				mockCache.EXPECT().
					IsUserInChat(gomock.Any(), chatID.String(), senderID.String()).
					Return(true, nil)

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
			mock: func(mockBroker *mock_repository.MockMessageBroker, mockCache *mock_repository.MockMessengerCache, mockRepo *mock_repository.MockMessengerRepository) {
			},
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
			mock: func(mockBroker *mock_repository.MockMessageBroker, mockCache *mock_repository.MockMessengerCache, mockRepo *mock_repository.MockMessengerRepository) {
			},
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
			mock: func(mockBroker *mock_repository.MockMessageBroker, mockCache *mock_repository.MockMessengerCache, mockRepo *mock_repository.MockMessengerRepository) {
			},
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
			mock: func(mockBroker *mock_repository.MockMessageBroker, mockCache *mock_repository.MockMessengerCache, mockRepo *mock_repository.MockMessengerRepository) {
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc, mockRepo, mockBroker, mockCache := setupMockService(ctrl)

			tt.mock(mockBroker, mockCache, mockRepo)

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
