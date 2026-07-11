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
	"github.com/stretchr/testify/assert"
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
		mock              func(mock *mock_repository.MockMessengerRepository)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name: "Success",
			input: &entity.DirectChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 2},
			},
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					CreateDirectChat(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Error - Repository Failed",
			input: &entity.DirectChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 2},
			},
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					CreateDirectChat(gomock.Any(), gomock.Any()).
					Return(errRepoQueryError)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - Validation (Missing ParticipantIDs)",
			input: &entity.DirectChat{
				ChatID:         1,
				ParticipantIDs: nil,
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Invalid ParticipantIDs count - 1)",
			input: &entity.DirectChat{
				ChatID:         1,
				ParticipantIDs: []int64{1},
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Invalid ParticipantIDs count - 3)",
			input: &entity.DirectChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 2, 3},
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Invalid Participant ID - 0)",
			input: &entity.DirectChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 0},
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Invalid Participant ID - negative)",
			input: &entity.DirectChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, -5},
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Invalid ChatID - negative)",
			input: &entity.DirectChat{
				ChatID:         -1,
				ParticipantIDs: []int64{1, 2},
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc, mockRepo, _ := setupMockService(ctrl)

			tt.mock(mockRepo)

			err := svc.CreateDirectChat(context.Background(), tt.input)

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

func TestMessengerService_CreateGroupChat(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name              string
		input             *entity.GroupChat
		mock              func(mock *mock_repository.MockMessengerRepository)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name: "Success",
			input: &entity.GroupChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 2, 3},
				Name:           "Dev Team",
				OwnerID:        1,
			},
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					CreateGroupChat(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, m *model.GroupChat) error {
						m.ChatID = 100
						m.CreatedAt = now
						return nil
					})
			},
			wantErr: false,
		},
		{
			name: "Success - Boundary Values",
			input: &entity.GroupChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 2},
				Name:           strings.Repeat("A", 64),
				Title:          ptr(strings.Repeat("B", 64)),
				Description:    ptr(strings.Repeat("C", 255)),
				OwnerID:        1,
			},
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					CreateGroupChat(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, m *model.GroupChat) error {
						m.ChatID = 100
						m.CreatedAt = now
						return nil
					})
			},
			wantErr: false,
		},
		{
			name: "Error - Repository Failed",
			input: &entity.GroupChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 2, 3},
				Name:           "Dev Team",
				OwnerID:        1,
			},
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					CreateGroupChat(gomock.Any(), gomock.Any()).
					Return(errRepoQueryError)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - Validation (Missing Name - empty)",
			input: &entity.GroupChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 2},
				Name:           "",
				OwnerID:        1,
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Name too long - 65)",
			input: &entity.GroupChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 2},
				Name:           strings.Repeat("A", 65),
				OwnerID:        1,
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Title too long - 65)",
			input: &entity.GroupChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 2},
				Name:           "Dev Team",
				Title:          ptr(strings.Repeat("A", 65)),
				OwnerID:        1,
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Description too long - 256)",
			input: &entity.GroupChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 2},
				Name:           "Dev Team",
				Description:    ptr(strings.Repeat("A", 256)),
				OwnerID:        1,
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Missing OwnerID - 0)",
			input: &entity.GroupChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 2},
				Name:           "Dev Team",
				OwnerID:        0,
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Invalid OwnerID - negative)",
			input: &entity.GroupChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 2},
				Name:           "Dev Team",
				OwnerID:        -1,
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Less than 2 participants - 1)",
			input: &entity.GroupChat{
				ChatID:         1,
				ParticipantIDs: []int64{1},
				Name:           "Dev Team",
				OwnerID:        1,
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Empty participants - 0)",
			input: &entity.GroupChat{
				ChatID:         1,
				ParticipantIDs: []int64{},
				Name:           "Dev Team",
				OwnerID:        1,
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Invalid participant ID - negative)",
			input: &entity.GroupChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, -5},
				Name:           "Dev Team",
				OwnerID:        1,
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc, mockRepo, _ := setupMockService(ctrl)

			tt.mock(mockRepo)

			err := svc.CreateGroupChat(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				var appErr *errs.AppError
				if errors.As(err, &appErr) {
					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int64(100), tt.input.ChatID)
				assert.Equal(t, now, tt.input.CreatedAt)
			}
		})
	}
}

func TestMessengerService_GetChatsByUserID(t *testing.T) {
	tests := []struct {
		name              string
		userID            int64
		limit             int
		offset            int
		mock              func(mock *mock_repository.MockMessengerRepository)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name:   "Success",
			userID: 1,
			limit:  10,
			offset: 0,
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					GetChatsByUserID(gomock.Any(), int64(1), 10, 0).
					Return([]model.Chat{
						{ChatID: 1, ChatType: "direct", Name: ptr("Chat 1"), OwnerID: ptr(int64(1))},
					}, nil)
			},
			wantErr: false,
		},
		{
			name:   "Success - Empty Result",
			userID: 2,
			limit:  10,
			offset: 0,
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					GetChatsByUserID(gomock.Any(), int64(2), 10, 0).
					Return([]model.Chat{}, nil)
			},
			wantErr: false,
		},
		{
			name:   "Error - Repository Failed",
			userID: 1,
			limit:  10,
			offset: 0,
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					GetChatsByUserID(gomock.Any(), int64(1), 10, 0).
					Return(nil, errRepoQueryError)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name:   "Error - Repository Failed (Invalid UserID - negative)",
			userID: -1,
			limit:  10,
			offset: 0,
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					GetChatsByUserID(gomock.Any(), int64(-1), 10, 0).
					Return(nil, errRepoQueryError)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name:   "Error - Repository Failed (Invalid limit - negative)",
			userID: 1,
			limit:  -10,
			offset: 0,
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					GetChatsByUserID(gomock.Any(), int64(1), -10, 0).
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

			svc, mockRepo, _ := setupMockService(ctrl)

			tt.mock(mockRepo)

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
			}
		})
	}
}

func TestMessengerService_UpdateGroupChat(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name              string
		input             *entity.UpdateGroupChat
		mock              func(mock *mock_repository.MockMessengerRepository)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name: "Success",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: 1,
				ChatID:        1,
				Name:          ptr("New Name"),
			},
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					UpdateGroupChat(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, m *model.UpdateGroupChat) error {
						m.CreatedAt = now
						return nil
					})
			},
			wantErr: false,
		},
		{
			name: "Success - Boundary Values",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: 1,
				ChatID:        1,
				Name:          ptr(strings.Repeat("A", 64)),
				Title:         ptr(strings.Repeat("B", 64)),
				Description:   ptr(strings.Repeat("C", 255)),
				OwnerID:       ptr(int64(1)),
			},
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					UpdateGroupChat(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, m *model.UpdateGroupChat) error {
						m.CreatedAt = now
						return nil
					})
			},
			wantErr: false,
		},
		{
			name: "Error - Repository Failed",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: 1,
				ChatID:        1,
				Name:          ptr("New Name"),
			},
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					UpdateGroupChat(gomock.Any(), gomock.Any()).
					Return(errRepoQueryError)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - Validation (Missing UserIDUpdater - 0)",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: 0,
				ChatID:        1,
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Invalid UserIDUpdater - negative)",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: -1,
				ChatID:        1,
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Missing ChatID - 0)",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: 1,
				ChatID:        0,
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Invalid ChatID - negative)",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: 1,
				ChatID:        -1,
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Invalid OwnerID - 0)",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: 1,
				ChatID:        1,
				OwnerID:       ptr(int64(0)),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Invalid OwnerID - negative)",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: 1,
				ChatID:        1,
				OwnerID:       ptr(int64(-5)),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Name too long - 65)",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: 1,
				ChatID:        1,
				Name:          ptr(strings.Repeat("A", 65)),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Name empty string)",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: 1,
				ChatID:        1,
				Name:          ptr(""),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Title too long - 65)",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: 1,
				ChatID:        1,
				Title:         ptr(strings.Repeat("A", 65)),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Description too long - 256)",
			input: &entity.UpdateGroupChat{
				UserIDUpdater: 1,
				ChatID:        1,
				Description:   ptr(strings.Repeat("A", 256)),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc, mockRepo, _ := setupMockService(ctrl)

			tt.mock(mockRepo)

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
		mock              func(mock *mock_repository.MockMessageBroker)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name: "Success",
			input: &entity.SendMessage{
				ChatID:   1,
				SenderID: 2,
				Message:  "Hello World",
			},
			mock: func(mock *mock_repository.MockMessageBroker) {
				mock.EXPECT().
					SendMessage(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Success - Boundary Values",
			input: &entity.SendMessage{
				ChatID:   1,
				SenderID: 2,
				Message:  strings.Repeat("A", 4096),
			},
			mock: func(mock *mock_repository.MockMessageBroker) {
				mock.EXPECT().
					SendMessage(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Error - Broker Failed",
			input: &entity.SendMessage{
				ChatID:   1,
				SenderID: 2,
				Message:  "Hello World",
			},
			mock: func(mock *mock_repository.MockMessageBroker) {
				mock.EXPECT().
					SendMessage(gomock.Any(), gomock.Any()).
					Return(errBrokerError)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - Validation (Missing ChatID - 0)",
			input: &entity.SendMessage{
				ChatID:   0,
				SenderID: 2,
				Message:  "Hello World",
			},
			mock:              func(mock *mock_repository.MockMessageBroker) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Invalid ChatID - negative)",
			input: &entity.SendMessage{
				ChatID:   -1,
				SenderID: 2,
				Message:  "Hello World",
			},
			mock:              func(mock *mock_repository.MockMessageBroker) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Missing SenderID - 0)",
			input: &entity.SendMessage{
				ChatID:   1,
				SenderID: 0,
				Message:  "Hello World",
			},
			mock:              func(mock *mock_repository.MockMessageBroker) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Invalid SenderID - negative)",
			input: &entity.SendMessage{
				ChatID:   1,
				SenderID: -5,
				Message:  "Hello World",
			},
			mock:              func(mock *mock_repository.MockMessageBroker) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Empty Message)",
			input: &entity.SendMessage{
				ChatID:   1,
				SenderID: 2,
				Message:  "",
			},
			mock:              func(mock *mock_repository.MockMessageBroker) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation (Message too long - 4097)",
			input: &entity.SendMessage{
				ChatID:   1,
				SenderID: 2,
				Message:  strings.Repeat("A", 4097),
			},
			mock:              func(mock *mock_repository.MockMessageBroker) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc, _, mockBroker := setupMockService(ctrl)

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
