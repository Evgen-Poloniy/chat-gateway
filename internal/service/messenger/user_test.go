// package messenger_test

// import (
// 	"context"
// 	"errors"
// 	"testing"
// 	"time"

// 	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
// 	"github.com/Evgen-Poloniy/chat-gateway/internal/model"
// 	mock_repository "github.com/Evgen-Poloniy/chat-gateway/internal/service/messenger/mocks"
// 	errs "github.com/Evgen-Poloniy/chat-gateway/pkg/errors"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"go.uber.org/mock/gomock"
// )

// var (
// 	errRepoUserNotFound = errs.NewAppError(errs.CodeUserNotFound, "user not found", nil)
// )

// func TestMessengerService_GetUserDataByUsername(t *testing.T) {
// 	now := time.Now()
// 	username := "testuser"

// 	tests := []struct {
// 		name              string
// 		input             string
// 		mock              func(mock *mock_repository.MockMessengerRepository)
// 		wantErr           bool
// 		expectedErrorCode errs.ErrCode
// 	}{
// 		{
// 			name:  "Success",
// 			input: username,
// 			mock: func(mock *mock_repository.MockMessengerRepository) {
// 				mock.EXPECT().
// 					GetUserDataByUsername(gomock.Any(), username).
// 					Return(&model.User{
// 						UserID:    1,
// 						Username:  username,
// 						Email:     ptr("test@test.com"),
// 						FirstName: ptr("John"),
// 						LastName:  ptr("Doe"),
// 						BirthDate: &now,
// 						CreatedAt: now,
// 						Gender:    ptr("man"),
// 					}, nil)
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:  "Error - User Not Found",
// 			input: username,
// 			mock: func(mock *mock_repository.MockMessengerRepository) {
// 				mock.EXPECT().
// 					GetUserDataByUsername(gomock.Any(), username).
// 					Return(nil, errRepoUserNotFound)
// 			},
// 			wantErr:           true,
// 			expectedErrorCode: errs.CodeUserNotFound,
// 		},
// 		{
// 			name:  "Error - Query Failed",
// 			input: username,
// 			mock: func(mock *mock_repository.MockMessengerRepository) {
// 				mock.EXPECT().
// 					GetUserDataByUsername(gomock.Any(), username).
// 					Return(nil, errRepoQueryError)
// 			},
// 			wantErr:           true,
// 			expectedErrorCode: errs.CodeQueryError,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			svc, mockRepo, _ := setupMockService(ctrl)

// 			tt.mock(mockRepo)

// 			user, err := svc.GetUserDataByUsername(context.Background(), tt.input)

// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				var appErr *errs.AppError
// 				if errors.As(err, &appErr) {
// 					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
// 				}
// 				assert.Nil(t, user)
// 			} else {
// 				assert.NoError(t, err)
// 				require.NotNil(t, user)
// 				assert.Equal(t, tt.input, user.Username)
// 				assert.Equal(t, int64(1), user.UserID)
// 			}
// 		})
// 	}
// }

// func TestMessengerService_GetUserDataByUserID(t *testing.T) {
// 	now := time.Now()
// 	var userID int64 = 1

// 	tests := []struct {
// 		name              string
// 		input             int64
// 		mock              func(mock *mock_repository.MockMessengerRepository)
// 		wantErr           bool
// 		expectedErrorCode errs.ErrCode
// 	}{
// 		{
// 			name:  "Success",
// 			input: userID,
// 			mock: func(mock *mock_repository.MockMessengerRepository) {
// 				mock.EXPECT().
// 					GetUserDataByUserID(gomock.Any(), userID).
// 					Return(&model.User{
// 						UserID:    userID,
// 						Username:  "testuser",
// 						Email:     ptr("test@test.com"),
// 						FirstName: ptr("John"),
// 						LastName:  ptr("Doe"),
// 						BirthDate: &now,
// 						CreatedAt: now,
// 						Gender:    ptr("man"),
// 					}, nil)
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name:  "Error - User Not Found",
// 			input: userID,
// 			mock: func(mock *mock_repository.MockMessengerRepository) {
// 				mock.EXPECT().
// 					GetUserDataByUserID(gomock.Any(), userID).
// 					Return(nil, errRepoUserNotFound)
// 			},
// 			wantErr:           true,
// 			expectedErrorCode: errs.CodeUserNotFound,
// 		},
// 		{
// 			name:  "Error - Query Failed",
// 			input: userID,
// 			mock: func(mock *mock_repository.MockMessengerRepository) {
// 				mock.EXPECT().
// 					GetUserDataByUserID(gomock.Any(), userID).
// 					Return(nil, errRepoQueryError)
// 			},
// 			wantErr:           true,
// 			expectedErrorCode: errs.CodeQueryError,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			svc, mockRepo, _ := setupMockService(ctrl)

// 			tt.mock(mockRepo)

// 			user, err := svc.GetUserDataByUserID(context.Background(), tt.input)

// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				var appErr *errs.AppError
// 				if errors.As(err, &appErr) {
// 					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
// 				}
// 				assert.Nil(t, user)
// 			} else {
// 				assert.NoError(t, err)
// 				require.NotNil(t, user)
// 				assert.Equal(t, tt.input, user.UserID)
// 				assert.Equal(t, "testuser", user.Username)
// 			}
// 		})
// 	}
// }

// func TestMessengerService_CreateUser(t *testing.T) {
// 	now := time.Now()

// 	tests := []struct {
// 		name              string
// 		input             *entity.User
// 		mock              func(mock *mock_repository.MockMessengerRepository, input *entity.User)
// 		wantErr           bool
// 		expectedErrorCode errs.ErrCode
// 	}{
// 		{
// 			name: "Success",
// 			input: &entity.User{
// 				UserID:    1,
// 				Username:  "newuser",
// 				Email:     ptr("new@test.com"),
// 				FirstName: ptr("New"),
// 				LastName:  ptr("User"),
// 				BirthDate: &now,
// 				Gender:    ptr("woman"),
// 			},
// 			mock: func(mock *mock_repository.MockMessengerRepository, input *entity.User) {
// 				mock.EXPECT().
// 					CreateUser(gomock.Any(), gomock.Any()).
// 					DoAndReturn(func(ctx context.Context, u *model.User) error {
// 						u.UserID = 1
// 						u.CreatedAt = now
// 						return nil
// 					})
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Error - Validation Failed (Empty Username)",
// 			input: &entity.User{
// 				Username: "",
// 			},
// 			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.User) {},
// 			wantErr:           true,
// 			expectedErrorCode: errs.CodeValidationError,
// 		},
// 		{
// 			name: "Error - Validation Failed (Invalid Email)",
// 			input: &entity.User{
// 				Username: "validuser",
// 				Email:    ptr("invalid-email"),
// 			},
// 			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.User) {},
// 			wantErr:           true,
// 			expectedErrorCode: errs.CodeValidationError,
// 		},
// 		{
// 			name: "Error - Validation Failed (Invalid Gender)",
// 			input: &entity.User{
// 				Username: "validuser",
// 				Gender:   ptr("unknown"),
// 			},
// 			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.User) {},
// 			wantErr:           true,
// 			expectedErrorCode: errs.CodeValidationError,
// 		},
// 		{
// 			name: "Error - Repository Failed",
// 			input: &entity.User{
// 				UserID:   1,
// 				Username: "newuser",
// 				Email:    ptr("new@test.com"),
// 				Gender:   ptr("man"),
// 			},
// 			mock: func(mock *mock_repository.MockMessengerRepository, input *entity.User) {
// 				mock.EXPECT().
// 					CreateUser(gomock.Any(), gomock.Any()).
// 					Return(errRepoQueryError)
// 			},
// 			wantErr:           true,
// 			expectedErrorCode: errs.CodeQueryError,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			svc, mockRepo, _ := setupMockService(ctrl)

// 			tt.mock(mockRepo, tt.input)

// 			err := svc.CreateUser(context.Background(), tt.input)

// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				var appErr *errs.AppError
// 				if errors.As(err, &appErr) {
// 					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
// 				}
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, int64(1), tt.input.UserID)
// 				assert.False(t, tt.input.CreatedAt.IsZero())
// 			}
// 		})
// 	}
// }

// func TestMessengerService_UpdateUser(t *testing.T) {
// 	now := time.Now()

// 	tests := []struct {
// 		name              string
// 		input             *entity.UpdateUser
// 		mock              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser)
// 		wantErr           bool
// 		expectedErrorCode errs.ErrCode
// 	}{
// 		{
// 			name: "Success",
// 			input: &entity.UpdateUser{
// 				UserID:    1,
// 				Username:  ptr("updateduser"),
// 				Email:     ptr("updated@test.com"),
// 				FirstName: ptr("Update"),
// 				LastName:  ptr("User"),
// 				BirthDate: &now,
// 				Gender:    ptr("man"),
// 			},
// 			mock: func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {
// 				mock.EXPECT().
// 					UpdateUser(gomock.Any(), gomock.Any()).
// 					DoAndReturn(func(ctx context.Context, u *model.UpdateUser) error {
// 						u.Username = ptr("updated_in_db")
// 						return nil
// 					})
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Error - Validation Failed (Invalid UserID)",
// 			input: &entity.UpdateUser{
// 				UserID: 0,
// 			},
// 			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
// 			wantErr:           true,
// 			expectedErrorCode: errs.CodeValidationError,
// 		},
// 		{
// 			name: "Error - Validation Failed (Short Username)",
// 			input: &entity.UpdateUser{
// 				UserID:   1,
// 				Username: ptr("ab"),
// 			},
// 			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
// 			wantErr:           true,
// 			expectedErrorCode: errs.CodeValidationError,
// 		},
// 		{
// 			name: "Error - Validation Failed (Invalid Gender)",
// 			input: &entity.UpdateUser{
// 				UserID: 1,
// 				Gender: ptr("alien"),
// 			},
// 			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
// 			wantErr:           true,
// 			expectedErrorCode: errs.CodeValidationError,
// 		},
// 		{
// 			name: "Error - Repository Failed",
// 			input: &entity.UpdateUser{
// 				UserID:   1,
// 				Username: ptr("updateduser"),
// 				Gender:   ptr("woman"),
// 			},
// 			mock: func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {
// 				mock.EXPECT().
// 					UpdateUser(gomock.Any(), gomock.Any()).
// 					Return(errRepoQueryError)
// 			},
// 			wantErr:           true,
// 			expectedErrorCode: errs.CodeQueryError,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			svc, mockRepo, _ := setupMockService(ctrl)

// 			tt.mock(mockRepo, tt.input)

// 			err := svc.UpdateUser(context.Background(), tt.input)

// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				var appErr *errs.AppError
// 				if errors.As(err, &appErr) {
// 					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
// 				}
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, "updated_in_db", *tt.input.Username)
// 			}
// 		})
// 	}
// }

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
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	errRepoUserNotFound = errs.NewAppError(errs.CodeUserNotFound, "user not found", nil)
)

func TestMessengerService_GetUserDataByUsername(t *testing.T) {
	now := time.Now()
	username := "testuser"

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
						UserID:    1,
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
		{
			name:  "Error - Empty Username",
			input: "",
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					GetUserDataByUsername(gomock.Any(), "").
					Return(nil, errRepoUserNotFound)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc, mockRepo, _ := setupMockService(ctrl)

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
				assert.Equal(t, int64(1), user.UserID)
			}
		})
	}
}

func TestMessengerService_GetUserDataByUserID(t *testing.T) {
	now := time.Now()
	var userID int64 = 1

	tests := []struct {
		name              string
		input             int64
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
		{
			name:  "Error - Invalid UserID (Negative)",
			input: -1,
			mock: func(mock *mock_repository.MockMessengerRepository) {
				mock.EXPECT().
					GetUserDataByUserID(gomock.Any(), int64(-1)).
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
				UserID:    1,
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
						u.UserID = 1
						u.CreatedAt = now
						return nil
					})
			},
			wantErr: false,
		},
		{
			name: "Success - Boundary Values",
			input: &entity.User{
				UserID:    1,
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
						u.UserID = 1
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
				UserID:   1,
				Username: "ab",
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.User) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Username too long)",
			input: &entity.User{
				UserID:   1,
				Username: strings.Repeat("a", 33),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.User) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Username not alphanum)",
			input: &entity.User{
				UserID:   1,
				Username: "user_name!",
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.User) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Invalid Email format)",
			input: &entity.User{
				UserID:   1,
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
				UserID:   1,
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
				UserID:    1,
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
				UserID:   1,
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
				UserID:   1,
				Username: "validuser",
				Gender:   ptr("unknown"),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.User) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (UserID 0)",
			input: &entity.User{
				UserID:   0,
				Username: "validuser",
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.User) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (UserID negative)",
			input: &entity.User{
				UserID:   -1,
				Username: "validuser",
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.User) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Repository Failed",
			input: &entity.User{
				UserID:   1,
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

			svc, mockRepo, _ := setupMockService(ctrl)

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
				assert.Equal(t, int64(1), tt.input.UserID)
				assert.False(t, tt.input.CreatedAt.IsZero())
			}
		})
	}
}

func TestMessengerService_UpdateUser(t *testing.T) {
	now := time.Now()

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
				UserID:    1,
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
				UserID:    1,
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
			name: "Error - Validation Failed (Invalid UserID - 0)",
			input: &entity.UpdateUser{
				UserID: 0,
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Invalid UserID - negative)",
			input: &entity.UpdateUser{
				UserID: -1,
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Username too short)",
			input: &entity.UpdateUser{
				UserID:   1,
				Username: ptr("ab"),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Username too long)",
			input: &entity.UpdateUser{
				UserID:   1,
				Username: ptr(strings.Repeat("a", 33)),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Username not alphanum)",
			input: &entity.UpdateUser{
				UserID:   1,
				Username: ptr("user@name"),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Username empty string)",
			input: &entity.UpdateUser{
				UserID:   1,
				Username: ptr(""),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Invalid Email format)",
			input: &entity.UpdateUser{
				UserID: 1,
				Email:  ptr("invalid-email"),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Email too long)",
			input: &entity.UpdateUser{
				UserID: 1,
				Email:  ptr(strings.Repeat("a", 60) + "@test.com"),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (FirstName too long)",
			input: &entity.UpdateUser{
				UserID:    1,
				FirstName: ptr(strings.Repeat("a", 101)),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (LastName too long)",
			input: &entity.UpdateUser{
				UserID:   1,
				LastName: ptr(strings.Repeat("a", 101)),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Validation Failed (Invalid Gender)",
			input: &entity.UpdateUser{
				UserID: 1,
				Gender: ptr("alien"),
			},
			mock:              func(mock *mock_repository.MockMessengerRepository, input *entity.UpdateUser) {},
			wantErr:           true,
			expectedErrorCode: errs.CodeValidationError,
		},
		{
			name: "Error - Repository Failed",
			input: &entity.UpdateUser{
				UserID:   1,
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

			svc, mockRepo, _ := setupMockService(ctrl)

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
