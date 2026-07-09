package pg_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
	pg "github.com/Evgen-Poloniy/chat-gateway/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresRepository_GetUserDataByUsername(t *testing.T) {
	now := time.Now()
	username := "testuser"

	query := `
        SELECT *
        FROM users
        WHERE username = $1`

	tests := []struct {
		name    string
		input   string
		mock    func(mock sqlmock.Sqlmock)
		wantErr bool
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
			name:  "Not Found",
			input: username,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(username).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
		{
			name:  "Query Error",
			input: username,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(username).
					WillReturnError(errors.New("db connection failed"))
			},
			wantErr: true,
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
        SELECT *
        FROM users
        WHERE id = $1`

	tests := []struct {
		name    string
		input   int64
		mock    func(mock sqlmock.Sqlmock)
		wantErr bool
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
			name:  "Not Found",
			input: userID,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(userID).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
		{
			name:  "Query Error",
			input: userID,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(userID).
					WillReturnError(errors.New("db query failed"))
			},
			wantErr: true,
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

	query := `INSERT INTO users (username, email, first_name, last_name, birth_date, gender) VALUES (?, ?, ?, ?, ?, ?) RETURNING id, created_at`

	tests := []struct {
		name    string
		input   *entity.User
		mock    func(mock sqlmock.Sqlmock, user *entity.User)
		wantErr bool
	}{
		{
			name: "Success",
			input: &entity.User{
				Username:  "newuser",
				Email:     ptr("new@test.com"),
				FirstName: ptr("New"),
				LastName:  ptr("User"),
				BirthDate: &now,
				Gender:    ptr("female"),
			},
			mock: func(mock sqlmock.Sqlmock, user *entity.User) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(user.Username, *user.Email, *user.FirstName, *user.LastName, *user.BirthDate, *user.Gender).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, now))
			},
			wantErr: false,
		},
		{
			name: "Error - Unique Violation",
			input: &entity.User{
				Username:  "newuser",
				Email:     ptr("new@test.com"),
				FirstName: ptr("New"),
				LastName:  ptr("User"),
				BirthDate: &now,
				Gender:    ptr("female"),
			},
			mock: func(mock sqlmock.Sqlmock, user *entity.User) {
				pgErr := &pgconn.PgError{
					Code:    "23505",
					Message: "duplicate key value violates unique constraint",
				}
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(user.Username, *user.Email, *user.FirstName, *user.LastName, *user.BirthDate, *user.Gender).
					WillReturnError(pgErr)
			},
			wantErr: true,
		},
		{
			name: "Error - Query Failed",
			input: &entity.User{
				Username:  "newuser",
				Email:     ptr("new@test.com"),
				FirstName: ptr("New"),
				LastName:  ptr("User"),
				BirthDate: &now,
				Gender:    ptr("female"),
			},
			mock: func(mock sqlmock.Sqlmock, user *entity.User) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(user.Username, *user.Email, *user.FirstName, *user.LastName, *user.BirthDate, *user.Gender).
					WillReturnError(errors.New("generic db error"))
			},
			wantErr: true,
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

	expectedQuery := `UPDATE users SET username = COALESCE(?, username), email = COALESCE(?, email), first_name = COALESCE(?, first_name), last_name = COALESCE(?, last_name), birth_date = COALESCE(?, birth_date), gender = COALESCE(?, gender) WHERE id = ?`

	tests := []struct {
		name    string
		input   *entity.UpdateUser
		mock    func(mock sqlmock.Sqlmock, user *entity.UpdateUser)
		wantErr bool
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
				Gender:    ptr("male"),
			},
			mock: func(mock sqlmock.Sqlmock, user *entity.UpdateUser) {
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(
						*user.Username, *user.Email, *user.FirstName,
						*user.LastName, *user.BirthDate, *user.Gender, user.UserID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Error - Unique Violation",
			input: &entity.UpdateUser{
				UserID:    1,
				Username:  ptr("updateduser"),
				Email:     ptr("updated@test.com"),
				FirstName: ptr("Update"),
				LastName:  ptr("User"),
				BirthDate: &now,
				Gender:    ptr("male"),
			},
			mock: func(mock sqlmock.Sqlmock, user *entity.UpdateUser) {
				pgErr := &pgconn.PgError{
					Code:    "23505",
					Message: "duplicate key value violates unique constraint",
				}
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(
						*user.Username, *user.Email, *user.FirstName,
						*user.LastName, *user.BirthDate, *user.Gender, user.UserID,
					).
					WillReturnError(pgErr)
			},
			wantErr: true,
		},
		{
			name: "Error - Query Failed",
			input: &entity.UpdateUser{
				UserID:    1,
				Username:  ptr("updateduser"),
				Email:     ptr("updated@test.com"),
				FirstName: ptr("Update"),
				LastName:  ptr("User"),
				BirthDate: &now,
				Gender:    ptr("male"),
			},
			mock: func(mock sqlmock.Sqlmock, user *entity.UpdateUser) {
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(
						*user.Username, *user.Email, *user.FirstName,
						*user.LastName, *user.BirthDate, *user.Gender, user.UserID,
					).
					WillReturnError(errors.New("generic db error"))
			},
			wantErr: true,
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
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
