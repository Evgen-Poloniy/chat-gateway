package messenger_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
	"github.com/Evgen-Poloniy/chat-gateway/internal/model"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestMessengerService_CreateDirectChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc, mockRepo, _ := setupMockService(ctrl)
	ctx := context.Background()

	tests := []struct {
		name    string
		chat    *entity.DirectChat
		mock    func()
		wantErr bool
	}{
		{
			name: "Success",
			chat: &entity.DirectChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 2},
			},
			mock: func() {
				mockRepo.EXPECT().
					CreateDirectChat(gomock.Any(), &model.DirectChat{ParticipantIDs: []int64{1, 2}}).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Repository error",
			chat: &entity.DirectChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 2},
			},
			mock: func() {
				mockRepo.EXPECT().
					CreateDirectChat(gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "Validation: Missing ChatID",
			chat: &entity.DirectChat{
				ChatID:         0,
				ParticipantIDs: []int64{1, 2},
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Only 1 participant",
			chat: &entity.DirectChat{
				ChatID:         1,
				ParticipantIDs: []int64{1},
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: 3 participants",
			chat: &entity.DirectChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 2, 3},
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Invalid participant ID",
			chat: &entity.DirectChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 0},
			},
			mock:    func() {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}
			err := svc.CreateDirectChat(ctx, tt.chat)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMessengerService_CreateGroupChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc, mockRepo, _ := setupMockService(ctrl)
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name    string
		chat    *entity.GroupChat
		mock    func()
		wantErr bool
	}{
		{
			name: "Success",
			chat: &entity.GroupChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 2, 3},
				Name:           "Dev Team",
				OwnerID:        1,
			},
			mock: func() {
				mockRepo.EXPECT().
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
			name: "Repository error",
			chat: &entity.GroupChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 2, 3},
				Name:           "Dev Team",
				OwnerID:        1,
			},
			mock: func() {
				mockRepo.EXPECT().
					CreateGroupChat(gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "Validation: Missing ChatID",
			chat: &entity.GroupChat{
				ChatID:         0,
				ParticipantIDs: []int64{1, 2},
				Name:           "Dev Team",
				OwnerID:        1,
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Less than 2 participants",
			chat: &entity.GroupChat{
				ChatID:         1,
				ParticipantIDs: []int64{1},
				Name:           "Dev Team",
				OwnerID:        1,
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Invalid participant ID",
			chat: &entity.GroupChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, -5},
				Name:           "Dev Team",
				OwnerID:        1,
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Missing Name",
			chat: &entity.GroupChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 2},
				Name:           "",
				OwnerID:        1,
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Name too long",
			chat: &entity.GroupChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 2},
				Name:           strings.Repeat("A", 65),
				OwnerID:        1,
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Missing OwnerID",
			chat: &entity.GroupChat{
				ChatID:         1,
				ParticipantIDs: []int64{1, 2},
				Name:           "Dev Team",
				OwnerID:        0,
			},
			mock:    func() {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}
			err := svc.CreateGroupChat(ctx, tt.chat)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int64(100), tt.chat.ChatID)
				assert.Equal(t, now, tt.chat.CreatedAt)
			}
		})
	}
}

func TestMessengerService_GetChatsByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc, mockRepo, _ := setupMockService(ctrl)
	ctx := context.Background()

	tests := []struct {
		name          string
		userID        int64
		limit         int
		offset        int
		mock          func()
		expectedChats []entity.Chat
		wantErr       bool
	}{
		{
			name:   "Success: return chats",
			userID: 1,
			limit:  10,
			offset: 0,
			mock: func() {
				mockRepo.EXPECT().
					GetChatsByUserID(ctx, int64(1), 10, 0).
					Return([]model.Chat{
						{ChatID: 1, ChatType: "direct", Name: ptr("Chat 1"), OwnerID: ptr(int64(1))},
						{ChatID: 2, ChatType: "group", Name: ptr("Chat 2"), OwnerID: ptr(int64(2))},
					}, nil)
			},
			expectedChats: []entity.Chat{
				{ChatID: 1, ChatType: "direct", Name: ptr("Chat 1"), OwnerID: ptr(int64(1))},
				{ChatID: 2, ChatType: "group", Name: ptr("Chat 2"), OwnerID: ptr(int64(2))},
			},
			wantErr: false,
		},
		{
			name:   "Success: empty result",
			userID: 2,
			limit:  10,
			offset: 0,
			mock: func() {
				mockRepo.EXPECT().
					GetChatsByUserID(ctx, int64(2), 10, 0).
					Return([]model.Chat{}, nil)
			},
			expectedChats: []entity.Chat{},
			wantErr:       false,
		},
		{
			name:   "Repository error",
			userID: 1,
			limit:  10,
			offset: 0,
			mock: func() {
				mockRepo.EXPECT().
					GetChatsByUserID(ctx, int64(1), 10, 0).
					Return(nil, errors.New("db error"))
			},
			expectedChats: nil,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			chats, err := svc.GetChatsByUserID(ctx, tt.userID, tt.limit, tt.offset)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, chats)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedChats, chats)
			}
		})
	}
}

func TestMessengerService_UpdateChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc, mockRepo, _ := setupMockService(ctrl)
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name    string
		chat    *entity.UpdateChat
		mock    func()
		wantErr bool
	}{
		{
			name: "Success: Update Name",
			chat: &entity.UpdateChat{
				UserIDUpdater: 1,
				ChatID:        1,
				Name:          ptr("New Name"),
			},
			mock: func() {
				mockRepo.EXPECT().
					UpdateChat(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, m *model.UpdateChat) error {
						m.CreatedAt = now
						return nil
					})
			},
			wantErr: false,
		},
		{
			name: "Success: Update Title and Description",
			chat: &entity.UpdateChat{
				UserIDUpdater: 1,
				ChatID:        1,
				Title:         ptr("New Title"),
				Description:   ptr("New Description"),
			},
			mock: func() {
				mockRepo.EXPECT().
					UpdateChat(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, m *model.UpdateChat) error {
						m.CreatedAt = now
						return nil
					})
			},
			wantErr: false,
		},
		{
			name: "Repository error",
			chat: &entity.UpdateChat{
				UserIDUpdater: 1,
				ChatID:        1,
				Name:          ptr("New Name"),
			},
			mock: func() {
				mockRepo.EXPECT().
					UpdateChat(gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "Validation: Missing UserIDUpdater",
			chat: &entity.UpdateChat{
				UserIDUpdater: 0,
				ChatID:        1,
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Missing ChatID",
			chat: &entity.UpdateChat{
				UserIDUpdater: 1,
				ChatID:        0,
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Name too long",
			chat: &entity.UpdateChat{
				UserIDUpdater: 1,
				ChatID:        1,
				Name:          ptr(strings.Repeat("A", 65)),
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Invalid OwnerID",
			chat: &entity.UpdateChat{
				UserIDUpdater: 1,
				ChatID:        1,
				OwnerID:       ptr(int64(0)),
			},
			mock:    func() {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}
			err := svc.UpdateChat(ctx, tt.chat)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, now, tt.chat.CreatedAt)
			}
		})
	}
}

func TestMessengerService_SendMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc, _, mockBroker := setupMockService(ctrl)
	ctx := context.Background()

	tests := []struct {
		name    string
		message *entity.SendMessage
		mock    func()
		wantErr bool
	}{
		{
			name: "Success",
			message: &entity.SendMessage{
				ChatID:   1,
				SenderID: 2,
				Message:  "Hello World",
			},
			mock: func() {
				expectedModel := &model.SendMessage{
					ChatID:   1,
					SenderID: 2,
					Message:  "Hello World",
				}
				mockBroker.EXPECT().
					SendMessage(gomock.Any(), expectedModel).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Broker error",
			message: &entity.SendMessage{
				ChatID:   1,
				SenderID: 2,
				Message:  "Hello World",
			},
			mock: func() {
				mockBroker.EXPECT().
					SendMessage(gomock.Any(), gomock.Any()).
					Return(errors.New("broker connection lost"))
			},
			wantErr: true,
		},
		{
			name: "Validation: Missing ChatID",
			message: &entity.SendMessage{
				ChatID:   0,
				SenderID: 2,
				Message:  "Hello World",
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Missing SenderID",
			message: &entity.SendMessage{
				ChatID:   1,
				SenderID: 0,
				Message:  "Hello World",
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Empty Message",
			message: &entity.SendMessage{
				ChatID:   1,
				SenderID: 2,
				Message:  "",
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Message too long",
			message: &entity.SendMessage{
				ChatID:   1,
				SenderID: 2,
				Message:  strings.Repeat("A", 4097),
			},
			mock:    func() {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}
			err := svc.SendMessage(ctx, tt.message)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
