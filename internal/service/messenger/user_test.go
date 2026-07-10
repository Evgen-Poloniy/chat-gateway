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

func TestMessengerService_GetUserDataByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc, mockRepo, _ := setupMockService(ctrl)
	ctx := context.Background()

	now := time.Now()
	username := "testuser"

	modelUser := &model.User{
		UserID:    1,
		Username:  username,
		Email:     ptr("test@example.com"),
		FirstName: ptr("John"),
		LastName:  ptr("Doe"),
		BirthDate: ptr(now),
		Gender:    ptr("man"),
		CreatedAt: now,
	}

	expectedUser := &entity.User{
		UserID:    1,
		Username:  username,
		Email:     ptr("test@example.com"),
		FirstName: ptr("John"),
		LastName:  ptr("Doe"),
		BirthDate: ptr(now),
		Gender:    ptr("man"),
		CreatedAt: now,
	}

	tests := []struct {
		name     string
		username string
		mock     func()
		want     *entity.User
		wantErr  bool
	}{
		{
			name:     "Success",
			username: username,
			mock: func() {
				mockRepo.EXPECT().
					GetUserDataByUsername(ctx, username).
					Return(modelUser, nil)
			},
			want:    expectedUser,
			wantErr: false,
		},
		{
			name:     "Repository error",
			username: username,
			mock: func() {
				mockRepo.EXPECT().
					GetUserDataByUsername(ctx, username).
					Return(nil, errors.New("not found"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := svc.GetUserDataByUsername(ctx, tt.username)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestMessengerService_GetUserDataByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc, mockRepo, _ := setupMockService(ctrl)
	ctx := context.Background()

	now := time.Now()
	userID := int64(1)

	modelUser := &model.User{
		UserID:    userID,
		Username:  "testuser",
		CreatedAt: now,
	}

	expectedUser := &entity.User{
		UserID:    userID,
		Username:  "testuser",
		CreatedAt: now,
	}

	tests := []struct {
		name    string
		userID  int64
		mock    func()
		want    *entity.User
		wantErr bool
	}{
		{
			name:   "Success",
			userID: userID,
			mock: func() {
				mockRepo.EXPECT().
					GetUserDataByUserID(ctx, userID).
					Return(modelUser, nil)
			},
			want:    expectedUser,
			wantErr: false,
		},
		{
			name:   "Repository error",
			userID: userID,
			mock: func() {
				mockRepo.EXPECT().
					GetUserDataByUserID(ctx, userID).
					Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := svc.GetUserDataByUserID(ctx, tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestMessengerService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc, mockRepo, _ := setupMockService(ctrl)
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name    string
		user    *entity.User
		mock    func()
		wantErr bool
	}{
		{
			name: "Success",
			user: &entity.User{
				UserID:   1,
				Username: "validUser123",
				Email:    ptr("test@test.com"),
				Gender:   ptr("man"),
			},
			mock: func() {
				mockRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, m *model.User) error {
						m.UserID = 100
						m.CreatedAt = now
						return nil
					})
			},
			wantErr: false,
		},
		{
			name: "Repository error",
			user: &entity.User{
				UserID:   1,
				Username: "validUser",
			},
			mock: func() {
				mockRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "Validation: Missing UserID",
			user: &entity.User{
				UserID:   0,
				Username: "validUser",
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Username too short",
			user: &entity.User{
				UserID:   1,
				Username: "ab",
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Username too long",
			user: &entity.User{
				UserID:   1,
				Username: strings.Repeat("a", 33),
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Username not alphanum",
			user: &entity.User{
				UserID:   1,
				Username: "user_name!",
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Invalid email format",
			user: &entity.User{
				UserID:   1,
				Username: "validUser",
				Email:    ptr("invalid-email"),
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Invalid gender",
			user: &entity.User{
				UserID:   1,
				Username: "validUser",
				Gender:   ptr("alien"),
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
			err := svc.CreateUser(ctx, tt.user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int64(100), tt.user.UserID)
				assert.Equal(t, now, tt.user.CreatedAt)
			}
		})
	}
}

func TestMessengerService_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc, mockRepo, _ := setupMockService(ctrl)
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name    string
		user    *entity.UpdateUser
		mock    func()
		wantErr bool
	}{
		{
			name: "Success: Update All Fields",
			user: &entity.UpdateUser{
				UserID:    1,
				Username:  ptr("newUser123"),
				Email:     ptr("new@email.com"),
				FirstName: ptr("Bob"),
				LastName:  ptr("Marley"),
				Gender:    ptr("woman"),
			},
			mock: func() {
				mockRepo.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, m *model.UpdateUser) error {
						m.CreatedAt = now
						return nil
					})
			},
			wantErr: false,
		},
		{
			name: "Repository error",
			user: &entity.UpdateUser{
				UserID: 1,
			},
			mock: func() {
				mockRepo.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Return(errors.New("db update error"))
			},
			wantErr: true,
		},
		{
			name: "Validation: Missing UserID",
			user: &entity.UpdateUser{
				UserID: 0,
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Username too short",
			user: &entity.UpdateUser{
				UserID:   1,
				Username: ptr("jo"),
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Username not alphanum",
			user: &entity.UpdateUser{
				UserID:   1,
				Username: ptr("h@cker"),
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Invalid email",
			user: &entity.UpdateUser{
				UserID: 1,
				Email:  ptr("not-an-email.com"),
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: Invalid gender",
			user: &entity.UpdateUser{
				UserID: 1,
				Gender: ptr("helicopter"),
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Validation: FirstName too long",
			user: &entity.UpdateUser{
				UserID:    1,
				FirstName: ptr(strings.Repeat("A", 101)),
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
			err := svc.UpdateUser(ctx, tt.user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, now, tt.user.CreatedAt)
			}
		})
	}
}
