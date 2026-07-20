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
	errRepoUserNotFound = errs.NewAppError(errs.CodeUserNotFound, "user not found", nil)
)

func TestMessengerService_GetUserDataByUsername(t *testing.T) {
	now := time.Now()
	username := "testuser"
	userID := uuid.New()

	tests := []struct {
		name              string
		input             string
		mock              func(mock *mock_repository.MockMessengerRepository)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name:  "Success",
			input: username,
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					GetUserDataByUsername(gomock.Any(), username).
					Return(&model.User{
						UserID:    userID,
						Username:  username,
						Email:     ptr("test@test.com"),
						FirstName: ptr("John"),
						LastName:  ptr("Doe"),
						BirthDate: &now,
						CreatedAt: now,
						Gender:    ptr("man"),
					}, nil)
			},
			wantErr: false,
		},
		{
			name:  "Error - User Not Found",
			input: username,
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					GetUserDataByUsername(gomock.Any(), username).
					Return(nil, errRepoUserNotFound)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeUserNotFound,
		},
		{
			name:  "Error - Query Failed",
			input: username,
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					GetUserDataByUsername(gomock.Any(), username).
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

			svc, mockRepo, _, _ := setupMockService(ctrl)

			tt.mock(mockRepo)

			user, err := svc.GetUserDataByUsername(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				var appErr *errs.AppError
				if errors.As(err, &appErr) {
					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
				}
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, tt.input, user.Username)
				assert.Equal(t, userID, user.UserID)
			}
		})
	}
}

func TestMessengerService_GetUserDataByUserID(t *testing.T) {
	now := time.Now()
	userID := uuid.New()

	tests := []struct {
		name              string
		input             uuid.UUID
		mock              func(mock *mock_repository.MockMessengerRepository)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name:  "Success",
			input: userID,
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					GetUserDataByUserID(gomock.Any(), userID).
					Return(&model.User{
						UserID:    userID,
						Username:  "testuser",
						Email:     ptr("test@test.com"),
						FirstName: ptr("John"),
						LastName:  ptr("Doe"),
						BirthDate: &now,
						CreatedAt: now,
						Gender:    ptr("man"),
					}, nil)
			},
			wantErr: false,
		},
		{
			name:  "Error - User Not Found",
			input: userID,
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					GetUserDataByUserID(gomock.Any(), userID).
					Return(nil, errRepoUserNotFound)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeUserNotFound,
		},
		{
			name:  "Error - Query Failed",
			input: userID,
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					GetUserDataByUserID(gomock.Any(), userID).
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

			svc, mockRepo, _, _ := setupMockService(ctrl)

			tt.mock(mockRepo)

			user, err := svc.GetUserDataByUserID(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				var appErr *errs.AppError
				if errors.As(err, &appErr) {
					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
				}
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, tt.input, user.UserID)
				assert.Equal(t, "testuser", user.Username)
			}
		})
	}
}

func TestMessengerService_CreateUser(t *testing.T) {
	now := time.Now()
	generatedUUID := uuid.New()

	tests := []struct {
		name              string
		input             *entity.User
		mock              func(mock *mock_repository.MockMessengerRepository, input *entity.User)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name: "Success",
			input: &entity.User{
				Username:  "newuser",
				Email:     ptr("new@test.com"),
				FirstName: ptr("New"),
				LastName:  ptr("User"),
				BirthDate: &now,
				Gender:    ptr("woman"),
			},
			mock: func(mock *mock_repository.MockMessengerRepository, input *entity.User) {
				mock.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, u *model.User) error {
						u.UserID = generatedUUID
						u.CreatedAt = now
						return nil
					})
			},
			wantErr: false,
		},
		{
			name: "Success - Boundary Values",
			input: &entity.User{
				Username:  strings.Repeat("a", 32),
				Email:     ptr(strings.Repeat("a", 50) + "@t.com"),
				FirstName: ptr(strings.Repeat("a", 100)),
				LastName:  ptr(strings.Repeat("a", 100)),
				BirthDate: &now,
				Gender:    ptr("man"),
			},
			mock: func(mock *mock_repository.MockMessengerRepository, input *entity.User) {
				mock.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, u *model.User) error {
						u.UserID = generatedUUID
						u.CreatedAt = now
						return nil
					})
			},
			wantErr: false,
		},
		{
			name: "Error - Validation Failed (Empty Username)",
			input: &entity.User{
				Username: "",
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.User) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Username too short)",
			input: &entity.User{
				Username: "ab",
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.User) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Username too long)",
			input: &entity.User{
				Username: strings.Repeat("a", 33),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.User) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Username not alphanum)",
			input: &entity.User{
				Username: "user_name!",
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.User) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Invalid Email format)",
			input: &entity.User{
				Username: "validuser",
				Email:    ptr("invalid-email"),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.User) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Email too long)",
			input: &entity.User{
				Username: "validuser",
				Email:    ptr(strings.Repeat("a", 60) + "@test.com"),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.User) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (FirstName too long)",
			input: &entity.User{
				Username:  "validuser",
				FirstName: ptr(strings.Repeat("a", 101)),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.User) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (LastName too long)",
			input: &entity.User{
				Username: "validuser",
				LastName: ptr(strings.Repeat("a", 101)),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.User) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Invalid Gender)",
			input: &entity.User{
				Username: "validuser",
				Gender:   ptr("unknown"),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.User) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Repository Failed",
			input: &entity.User{
				Username: "newuser",
				Email:    ptr("new@test.com"),
				Gender:   ptr("man"),
			},
			mock: func(mock *mock_repository.MockMessengerRepository, input *entity.User) {
				mock.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(errRepoQueryError)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc, mockRepo, _, _ := setupMockService(ctrl)

			tt.mock(mockRepo, tt.input)

			err := svc.CreateUser(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				var appErr *errs.AppError
				if errors.As(err, &appErr) {
					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, generatedUUID, tt.input.UserID)
				assert.False(t, tt.input.CreatedAt.IsZero())
			}
		})
	}
}

func TestMessengerService_GetUserIDsByChatID(t *testing.T) {
	chatID := uuid.New()
	userID1 := uuid.New()
	userID2 := uuid.New()
	userID3 := uuid.New()

	tests := []struct {
		name              string
		input             uuid.UUID
		mockCache         func(mock *mock_repository.MockMessengerCache)
		mockRepo          func(mock *mock_repository.MockMessengerRepository)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name:  "Success - From Cache",
			input: chatID,
			mockCache: func(mock *mock_repository.MockMessengerCache) {
				mock.EXPECT().
					GetChatMembers(gomock.Any(), chatID.String()).
					Return([]string{userID1.String(), userID2.String()}, nil)
			},
			mockRepo: func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:  false,
		},
		{
			name:  "Success - From DB (Cache Empty)",
			input: chatID,
			mockCache: func(mock *mock_repository.MockMessengerCache) {
				mock.EXPECT().
					GetChatMembers(gomock.Any(), chatID.String()).
					Return([]string{}, nil)
			},
			mockRepo: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					GetUserIDsByChatID(gomock.Any(), chatID).
					Return([]uuid.UUID{userID1, userID2, userID3}, nil)
			},
			wantErr: false,
		},
		{
			name:  "Error - Chat Not Found (DB Empty)",
			input: chatID,
			mockCache: func(mock *mock_repository.MockMessengerCache) {
				mock.EXPECT().
					GetChatMembers(gomock.Any(), chatID.String()).
					Return([]string{}, nil)
			},
			mockRepo: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					GetUserIDsByChatID(gomock.Any(), chatID).
					Return([]uuid.UUID{}, nil)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeChatNotFound,
		},
		{
			name:  "Error - Cache Failed",
			input: chatID,
			mockCache: func(mock *mock_repository.MockMessengerCache) {
				mock.EXPECT().
					GetChatMembers(gomock.Any(), chatID.String()).
					Return(nil, errCacheFailed)
			},
			mockRepo:          func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeRedisError,
		},
		{
			name:  "Error - DB Query Failed",
			input: chatID,
			mockCache: func(mock *mock_repository.MockMessengerCache) {
				mock.EXPECT().
					GetChatMembers(gomock.Any(), chatID.String()).
					Return([]string{}, nil)
			},
			mockRepo: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					GetUserIDsByChatID(gomock.Any(), chatID).
					Return(nil, errRepoQueryError)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name:  "Error - Invalid UUID in Cache",
			input: chatID,
			mockCache: func(mock *mock_repository.MockMessengerCache) {
				mock.EXPECT().
					GetChatMembers(gomock.Any(), chatID.String()).
					Return([]string{"invalid-uuid"}, nil)
			},
			mockRepo:          func(mock *mock_repository.MockMessengerRepository) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc, mockRepo, _, mockCache := setupMockService(ctrl)

			tt.mockCache(mockCache)
			tt.mockRepo(mockRepo)

			userIDs, err := svc.GetUserIDsByChatID(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				var appErr *errs.AppError
				if assert.ErrorAs(t, err, &appErr) {
					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
				}
				assert.Nil(t, userIDs)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, userIDs)
				if tt.name == "Success - From Cache" {
					assert.Len(t, userIDs, 2)
					assert.ElementsMatch(t, []uuid.UUID{userID1, userID2}, userIDs)
				} else if tt.name == "Success - From DB (Cache Empty)" {
					assert.Len(t, userIDs, 3)
					assert.ElementsMatch(t, []uuid.UUID{userID1, userID2, userID3}, userIDs)
				}
			}
		})
	}
}

func TestMessengerService_UpdateUser(t *testing.T) {
	now := time.Now()
	userID := uuid.New()

	tests := []struct {
		name              string
		input             *entity.UpdateUser
		mock              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name: "Success",
			input: &entity.UpdateUser{
				UserID:    userID,
				Username:  ptr("updateduser"),
				Email:     ptr("updated@test.com"),
				FirstName: ptr("Update"),
				LastName:  ptr("User"),
				BirthDate: &now,
				Gender:    ptr("man"),
			},
			mock: func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {
				mock.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, u *model.UpdateUser) error {
						u.Username = ptr("updated_in_db")
						return nil
					})
			},
			wantErr: false,
		},
		{
			name: "Success - Boundary Values",
			input: &entity.UpdateUser{
				UserID:    userID,
				Username:  ptr(strings.Repeat("a", 32)),
				Email:     ptr(strings.Repeat("a", 50) + "@t.com"),
				FirstName: ptr(strings.Repeat("a", 100)),
				LastName:  ptr(strings.Repeat("a", 100)),
				BirthDate: &now,
				Gender:    ptr("woman"),
			},
			mock: func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {
				mock.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, u *model.UpdateUser) error {
						u.Username = ptr("updated_in_db")
						return nil
					})
			},
			wantErr: false,
		},
		{
			name: "Error - Validation Failed (Invalid UserID - Nil UUID)",
			input: &entity.UpdateUser{
				UserID: uuid.Nil,
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Username too short)",
			input: &entity.UpdateUser{
				UserID:   userID,
				Username: ptr("ab"),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Username too long)",
			input: &entity.UpdateUser{
				UserID:   userID,
				Username: ptr(strings.Repeat("a", 33)),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Username not alphanum)",
			input: &entity.UpdateUser{
				UserID:   userID,
				Username: ptr("user@name"),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Username empty string)",
			input: &entity.UpdateUser{
				UserID:   userID,
				Username: ptr(""),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Invalid Email format)",
			input: &entity.UpdateUser{
				UserID: userID,
				Email:  ptr("invalid-email"),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Email too long)",
			input: &entity.UpdateUser{
				UserID: userID,
				Email:  ptr(strings.Repeat("a", 60) + "@test.com"),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (FirstName too long)",
			input: &entity.UpdateUser{
				UserID:    userID,
				FirstName: ptr(strings.Repeat("a", 101)),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (LastName too long)",
			input: &entity.UpdateUser{
				UserID:   userID,
				LastName: ptr(strings.Repeat("a", 101)),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Invalid Gender)",
			input: &entity.UpdateUser{
				UserID: userID,
				Gender: ptr("alien"),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Repository Failed",
			input: &entity.UpdateUser{
				UserID:   userID,
				Username: ptr("updateduser"),
				Gender:   ptr("woman"),
			},
			mock: func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {
				mock.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Return(errRepoQueryError)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc, mockRepo, _, _ := setupMockService(ctrl)

			tt.mock(mockRepo, tt.input)

			err := svc.UpdateUser(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				var appErr *errs.AppError
				if errors.As(err, &appErr) {
					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "updated_in_db", *tt.input.Username)
			}
		})
	}
}
