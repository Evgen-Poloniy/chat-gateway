package pg_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Evgen-Poloniy/chat-gateway/internal/model"
	pg "github.com/Evgen-Poloniy/chat-gateway/internal/repository/postgres"
	errs "github.com/Evgen-Poloniy/chat-gateway/pkg/errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresRepository_GetUserDataByUsername(t *testing.T) {
	now := time.Now()
	username := "testuser"

	query := `
        SELECT id, username, email, first_name, last_name, birth_date, created_at, gender
        FROM users
        WHERE username = $1`

	tests := []struct {
		name              string
		input             string
		mock              func(mock sqlmock.Sqlmock)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name:  "Success",
			input: username,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "first_name", "last_name", "birth_date", "created_at", "gender"}).
					AddRow(1, username, "test@test.com", "John", "Doe", now, now, "male")

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(username).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:  "Error - Not Found",
			input: username,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(username).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeUserNotFound,
		},
		{
			name:  "Error - Query Failed",
			input: username,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(username).
					WillReturnError(errDBQueryFailed)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			repo := pg.NewPostgresRepository(db)
			tt.mock(mock)

			user, err := repo.GetUserDataByUsername(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				var appErr *errs.AppError
				if assert.ErrorAs(t, err, &appErr) {
					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
				}
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, tt.input, user.Username)
				assert.Equal(t, int64(1), user.UserID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresRepository_GetUserDataByUserID(t *testing.T) {
	now := time.Now()
	var userID int64 = 1

	query := `
        SELECT id, username, email, first_name, last_name, birth_date, created_at, gender
        FROM users
        WHERE id = $1`

	tests := []struct {
		name              string
		input             int64
		mock              func(mock sqlmock.Sqlmock)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name:  "Success",
			input: userID,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "first_name", "last_name", "birth_date", "created_at", "gender"}).
					AddRow(userID, "testuser", "test@test.com", "John", "Doe", now, now, "male")

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:  "Error - Not Found",
			input: userID,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(userID).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeUserNotFound,
		},
		{
			name:  "Error - Query Failed",
			input: userID,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(userID).
					WillReturnError(errDBQueryFailed)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			repo := pg.NewPostgresRepository(db)
			tt.mock(mock)

			user, err := repo.GetUserDataByUserID(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				var appErr *errs.AppError
				if assert.ErrorAs(t, err, &appErr) {
					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
				}
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, tt.input, user.UserID)
				assert.Equal(t, "testuser", user.Username)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresRepository_CreateUser(t *testing.T) {
	now := time.Now()

	query := `INSERT INTO users
		(username, email, first_name, last_name, birth_date, gender)
		VALUES (?, ?, ?, ?, ?, ?) RETURNING id, created_at`

	tests := []struct {
		name              string
		input             *model.User
		mock              func(mock sqlmock.Sqlmock, user *model.User)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name: "Success",
			input: &model.User{
				Username:  "newuser",
				Email:     ptr("new@test.com"),
				FirstName: ptr("New"),
				LastName:  ptr("User"),
				BirthDate: &now,
				Gender:    ptr("female"),
			},
			mock: func(mock sqlmock.Sqlmock, user *model.User) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(user.Username, user.Email, user.FirstName, user.LastName, user.BirthDate, user.Gender).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, now))
			},
			wantErr: false,
		},
		{
			name: "Error - Bind Named Params",
			input: &model.User{
				Username:  "newuser",
				Email:     ptr("new@test.com"),
				FirstName: ptr("New"),
				LastName:  ptr("User"),
				BirthDate: &now,
				Gender:    ptr("female"),
			},
			mock: func(mock sqlmock.Sqlmock, user *model.User) {
				mock.ExpectQuery("INSERT INTO users *").
					WithArgs(user.Username, user.Email, user.FirstName, user.LastName, user.BirthDate, user.Gender).
					WillReturnError(errDBBindError)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - Unique Violation",
			input: &model.User{
				Username:  "newuser",
				Email:     ptr("new@test.com"),
				FirstName: ptr("New"),
				LastName:  ptr("User"),
				BirthDate: &now,
				Gender:    ptr("female"),
			},
			mock: func(mock sqlmock.Sqlmock, user *model.User) {
				pgErr := &pgconn.PgError{
					Code:    "23505",
					Message: "duplicate key value violates unique constraint",
				}
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(user.Username, user.Email, user.FirstName, user.LastName, user.BirthDate, user.Gender).
					WillReturnError(pgErr)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeUniqueViolation,
		},
		{
			name: "Error - Query Failed",
			input: &model.User{
				Username:  "newuser",
				Email:     ptr("new@test.com"),
				FirstName: ptr("New"),
				LastName:  ptr("User"),
				BirthDate: &now,
				Gender:    ptr("female"),
			},
			mock: func(mock sqlmock.Sqlmock, user *model.User) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(user.Username, user.Email, user.FirstName, user.LastName, user.BirthDate, user.Gender).
					WillReturnError(errDBQueryFailed)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			repo := pg.NewPostgresRepository(db)
			tt.mock(mock, tt.input)

			err := repo.CreateUser(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				var appErr *errs.AppError
				if assert.ErrorAs(t, err, &appErr) {
					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int64(1), tt.input.UserID)
				assert.NotNil(t, tt.input.CreatedAt)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresRepository_UpdateUser(t *testing.T) {
	now := time.Now()

	expectedQuery := `UPDATE users SET
		username = COALESCE(?, username),
		email = COALESCE(?, email),
		first_name = COALESCE(?, first_name),
		last_name = COALESCE(?, last_name),
		birth_date = COALESCE(?, birth_date),
		gender = COALESCE(?, gender)
		WHERE id = ?`

	tests := []struct {
		name              string
		input             *model.UpdateUser
		mock              func(mock sqlmock.Sqlmock, user *model.UpdateUser)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name: "Success",
			input: &model.UpdateUser{
				UserID:    1,
				Username:  ptr("updateduser"),
				Email:     ptr("updated@test.com"),
				FirstName: ptr("Update"),
				LastName:  ptr("User"),
				BirthDate: &now,
				Gender:    ptr("male"),
			},
			mock: func(mock sqlmock.Sqlmock, user *model.UpdateUser) {
				rows := sqlmock.NewRows([]string{"username", "email", "first_name", "last_name", "birth_date", "gender", "created_at"}).
					AddRow(user.Username, user.Email, user.FirstName, user.LastName, user.BirthDate, user.Gender, now)

				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(
						user.Username, user.Email, user.FirstName,
						user.LastName, user.BirthDate, user.Gender, user.UserID,
					).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "Error - Unique Violation",
			input: &model.UpdateUser{
				UserID:    1,
				Username:  ptr("updateduser"),
				Email:     ptr("updated@test.com"),
				FirstName: ptr("Update"),
				LastName:  ptr("User"),
				BirthDate: &now,
				Gender:    ptr("male"),
			},
			mock: func(mock sqlmock.Sqlmock, user *model.UpdateUser) {
				pgErr := &pgconn.PgError{
					Code:    "23505",
					Message: "duplicate key value violates unique constraint",
				}
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(
						user.Username, user.Email, user.FirstName,
						user.LastName, user.BirthDate, user.Gender, user.UserID,
					).
					WillReturnError(pgErr)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeUniqueViolation,
		},
		{
			name: "Error - Query Failed",
			input: &model.UpdateUser{
				UserID:    1,
				Username:  ptr("updateduser"),
				Email:     ptr("updated@test.com"),
				FirstName: ptr("Update"),
				LastName:  ptr("User"),
				BirthDate: &now,
				Gender:    ptr("male"),
			},
			mock: func(mock sqlmock.Sqlmock, user *model.UpdateUser) {
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(
						user.Username, user.Email, user.FirstName,
						user.LastName, user.BirthDate, user.Gender, user.UserID,
					).
					WillReturnError(errDBQueryFailed)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - Bind Named Params",
			input: &model.UpdateUser{
				UserID:    1,
				Username:  ptr("updateduser"),
				Email:     ptr("updated@test.com"),
				FirstName: ptr("Update"),
				LastName:  ptr("User"),
				BirthDate: &now,
				Gender:    ptr("male"),
			},
			mock: func(mock sqlmock.Sqlmock, user *model.UpdateUser) {
				mock.ExpectQuery("UPDATE users SET *").
					WithArgs(
						user.Username, user.Email, user.FirstName,
						user.LastName, user.BirthDate, user.Gender, user.UserID,
					).
					WillReturnError(errDBBindError)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - Not Found",
			input: &model.UpdateUser{
				UserID:    1,
				Username:  ptr("updateduser"),
				Email:     ptr("updated@test.com"),
				FirstName: ptr("Update"),
				LastName:  ptr("User"),
				BirthDate: &now,
				Gender:    ptr("male"),
			},
			mock: func(mock sqlmock.Sqlmock, user *model.UpdateUser) {
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(
						user.Username, user.Email, user.FirstName,
						user.LastName, user.BirthDate, user.Gender, user.UserID,
					).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			repo := pg.NewPostgresRepository(db)
			tt.mock(mock, tt.input)

			err := repo.UpdateUser(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				var appErr *errs.AppError
				if assert.ErrorAs(t, err, &appErr) {
					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
				}
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
