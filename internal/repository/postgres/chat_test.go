package pg_test

import (
	"context"
	"database/sql"
	"errors"
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

var (
	errConnectionLost    = errors.New("connection lost")
	errSyntaxError       = errors.New("syntax error")
	errFKViolation       = errors.New("foreign key violation")
	errCommitFailed      = errors.New("commit failed")
	errMemberInsertError = errors.New("member insert error")
)

func TestPostgresRepository_CreateDirectChat(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name              string
		input             *model.DirectChat
		mock              func(mock sqlmock.Sqlmock)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name: "Success",
			input: &model.DirectChat{
				ParticipantIDs: []int64{1, 2},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (type) VALUES ('direct') RETURNING id, created_at`)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, now))

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO chat_members (chat_id, user_id) VALUES ($1, $2), ($1, $3)`)).
					WithArgs(1, 1, 2).
					WillReturnResult(sqlmock.NewResult(1, 2))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Error - Transaction Begin Failed",
			input: &model.DirectChat{
				ParticipantIDs: []int64{1, 2},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errConnectionLost)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeTransactionError,
		},
		{
			name: "Error - Unique Violation",
			input: &model.DirectChat{
				ParticipantIDs: []int64{1, 2},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				pgErr := &pgconn.PgError{
					Code:    "23505",
					Message: "duplicate key value violates unique constraint",
					Detail:  "Key already exists.",
				}

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (type) VALUES ('direct') RETURNING id, created_at`)).
					WillReturnError(pgErr)

				mock.ExpectRollback()
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeUniqueViolation,
		},
		{
			name: "Error - Query Failed",
			input: &model.DirectChat{
				ParticipantIDs: []int64{1, 2},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (type) VALUES ('direct') RETURNING id, created_at`)).
					WillReturnError(errSyntaxError)

				mock.ExpectRollback()
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - Insert Members Failed",
			input: &model.DirectChat{
				ParticipantIDs: []int64{1, 2},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (type) VALUES ('direct') RETURNING id, created_at`)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, now))

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO chat_members (chat_id, user_id) VALUES ($1, $2), ($1, $3)`)).
					WithArgs(1, 1, 2).
					WillReturnError(errFKViolation)

				mock.ExpectRollback()
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - Commit Failed",
			input: &model.DirectChat{
				ParticipantIDs: []int64{1, 2},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (type) VALUES ('direct') RETURNING id, created_at`)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, now))

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO chat_members (chat_id, user_id) VALUES ($1, $2), ($1, $3)`)).
					WithArgs(1, 1, 2).
					WillReturnResult(sqlmock.NewResult(1, 2))

				mock.ExpectCommit().WillReturnError(errCommitFailed)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeTransactionError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			repo := pg.NewPostgresRepository(db)

			tt.mock(mock)

			err := repo.CreateDirectChat(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				var appErr *errs.AppError
				if assert.ErrorAs(t, err, &appErr) {
					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int64(1), tt.input.ChatID)
				assert.Equal(t, now, tt.input.CreatedAt)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresRepository_CreateGroupChat(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name              string
		input             *model.GroupChat
		mock              func(mock sqlmock.Sqlmock)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name: "Success",
			input: &model.GroupChat{
				Name:           "Group 1",
				Title:          ptr("Group chat"),
				Description:    ptr("Group chat"),
				OwnerID:        1,
				ParticipantIDs: []int64{1, 2, 3},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (type, name, title, description, owner_id) VALUES ('group', $1, $2, $3, $4) RETURNING id, created_at`)).
					WithArgs("Group 1", "Group chat", "Group chat", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(2, now))

				expectedQuery := regexp.QuoteMeta(`INSERT INTO chat_members (chat_id, user_id) VALUES ($1, $2), ($1, $3), ($1, $4)`)
				mock.ExpectExec(expectedQuery).
					WithArgs(2, 1, 2, 3).
					WillReturnResult(sqlmock.NewResult(1, 3))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Error - Transaction Begin Failed",
			input: &model.GroupChat{
				Name: "Group 2",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(context.DeadlineExceeded)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeTransactionError,
		},
		{
			name: "Error - Unique Violation",
			input: &model.GroupChat{
				Name:           "Duplicate",
				Title:          nil,
				Description:    nil,
				OwnerID:        1,
				ParticipantIDs: []int64{1},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				pgErr := &pgconn.PgError{
					Code:    "23505",
					Message: "unique violation",
					Detail:  "detail",
				}

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (type, name, title, description, owner_id) VALUES ('group', $1, $2, $3, $4) RETURNING id, created_at`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(pgErr)

				mock.ExpectRollback()
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeUniqueViolation,
		},
		{
			name: "Error - Foreign Key Violation",
			input: &model.GroupChat{
				Name:           "New Group",
				OwnerID:        999,
				Title:          nil,
				Description:    nil,
				ParticipantIDs: []int64{999},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				pgErr := &pgconn.PgError{Code: "23503", Message: "fk violation", Detail: "owner_id does not exist"}
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (type, name, title, description, owner_id) VALUES ('group', $1, $2, $3, $4) RETURNING id, created_at`)).
					WithArgs("New Group", nil, nil, int64(999)).
					WillReturnError(pgErr)

				mock.ExpectRollback()
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeForeignKeyViolation,
		},
		{
			name: "Error - Insert Members Failed",
			input: &model.GroupChat{
				Name:           "Group",
				OwnerID:        1,
				Title:          nil,
				Description:    nil,
				ParticipantIDs: []int64{1, 2},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (type, name, title, description, owner_id) VALUES ('group', $1, $2, $3, $4) RETURNING id, created_at`)).
					WithArgs("Group", nil, nil, int64(1)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(10, now))

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO chat_members (chat_id, user_id) VALUES ($1, $2), ($1, $3)`)).
					WithArgs(int64(10), int64(1), int64(2)).
					WillReturnError(errMemberInsertError)

				mock.ExpectRollback()
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - Query Failed",
			input: &model.GroupChat{
				Name:           "New Group",
				OwnerID:        999,
				Title:          nil,
				Description:    nil,
				ParticipantIDs: []int64{999},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (type, name, title, description, owner_id) VALUES ('group', $1, $2, $3, $4) RETURNING id, created_at`)).
					WithArgs("New Group", nil, nil, int64(999)).
					WillReturnError(errDBQueryFailed)

				mock.ExpectRollback()
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - Commit Failed",
			input: &model.GroupChat{
				Name:           "Group",
				OwnerID:        1,
				Title:          nil,
				Description:    nil,
				ParticipantIDs: []int64{1},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (type, name, title, description, owner_id) VALUES ('group', $1, $2, $3, $4) RETURNING id, created_at`)).
					WithArgs("Group", nil, nil, int64(1)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(10, now))

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO chat_members (chat_id, user_id) VALUES ($1, $2)`)).
					WithArgs(int64(10), int64(1)).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit().WillReturnError(errCommitFailed)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeTransactionError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupMockDB(t)
			defer cleanup()

			repo := pg.NewPostgresRepository(db)

			tt.mock(mock)

			err := repo.CreateGroupChat(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				var appErr *errs.AppError
				if assert.ErrorAs(t, err, &appErr) {
					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int64(2), tt.input.ChatID)
				assert.Equal(t, now, tt.input.CreatedAt)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresRepository_GetChatsByUserID(t *testing.T) {
	now := time.Now()

	query := `
        SELECT
            c.id,
            c.type,
            c.name,
            c.title,
            c.description,
            c.created_at,
            c.owner_id
        FROM chat_members cm
        JOIN chats c ON cm.chat_id = c.id
        WHERE cm.user_id = $1
        LIMIT $2 OFFSET $3`

	tests := []struct {
		name              string
		userID            int64
		limit             int
		offset            int
		mock              func(mock sqlmock.Sqlmock)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name:   "Success",
			userID: 1,
			limit:  10,
			offset: 0,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "type", "name", "title", "description", "created_at", "owner_id"}).
					AddRow(int64(5), "group", "group", nil, nil, now, int64(1))

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(int64(1), 10, 0).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:   "Error - Chats Not Found",
			userID: 1,
			limit:  10,
			offset: 0,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "type", "name", "title", "description", "created_at", "owner_id"})

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(int64(1), 10, 0).
					WillReturnRows(rows)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeChatNotFound,
		},
		{
			name:   "Error - Query Failed",
			userID: 1,
			limit:  10,
			offset: 0,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(int64(1), 10, 0).
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

			chats, err := repo.GetChatsByUserID(context.Background(), tt.userID, tt.limit, tt.offset)

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
				require.Len(t, chats, 1)

				assert.Equal(t, int64(5), chats[0].ChatID)
				assert.Equal(t, "group", chats[0].ChatType)
				assert.Equal(t, "group", *chats[0].Name)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresRepository_UpdateChat(t *testing.T) {
	now := time.Now()

	expectedQuery := `
		UPDATE chats
		SET
			name = COALESCE(?, name),
			title = COALESCE(?, title),
			description = COALESCE(?, description),
			owner_id = CASE
				WHEN ? IS NULL THEN owner_id
				WHEN ? = ? THEN owner_id
				ELSE owner_id
			END
		WHERE id = ?
		RETURNING name, title, description, created_at, owner_id
	`

	tests := []struct {
		name              string
		input             *model.UpdateChat
		mock              func(mock sqlmock.Sqlmock, chat *model.UpdateChat)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name: "Success",
			input: &model.UpdateChat{
				ChatID:        10,
				Name:          ptr("Updated Name"),
				Title:         ptr("Updated Title"),
				Description:   ptr("Updated Description"),
				OwnerID:       ptr(int64(2)),
				UserIDUpdater: 1,
			},
			mock: func(mock sqlmock.Sqlmock, chat *model.UpdateChat) {
				rows := sqlmock.NewRows([]string{"name", "title", "description", "created_at", "owner_id"}).
					AddRow(*chat.Name, *chat.Title, *chat.Description, now, *chat.OwnerID)

				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(
						chat.Name,
						chat.Title,
						chat.Description,
						chat.OwnerID,
						chat.OwnerID,
						chat.UserIDUpdater,
						chat.ChatID,
					).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "Error - Chat Not Found",
			input: &model.UpdateChat{
				ChatID:        10,
				Name:          ptr("Updated Name"),
				Title:         ptr("Updated Title"),
				Description:   ptr("Updated Description"),
				OwnerID:       ptr(int64(2)),
				UserIDUpdater: 1,
			},
			mock: func(mock sqlmock.Sqlmock, chat *model.UpdateChat) {
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(
						chat.Name, chat.Title, chat.Description,
						chat.OwnerID, chat.OwnerID, chat.UserIDUpdater, chat.ChatID,
					).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeChatNotFound,
		},
		{
			name: "Error - Unique Violation",
			input: &model.UpdateChat{
				ChatID:        10,
				Name:          ptr("Updated Name"),
				Title:         ptr("Updated Title"),
				Description:   ptr("Updated Description"),
				OwnerID:       ptr(int64(2)),
				UserIDUpdater: 1,
			},
			mock: func(mock sqlmock.Sqlmock, chat *model.UpdateChat) {
				pgErr := &pgconn.PgError{
					Code:    "23505",
					Message: "duplicate key value violates unique constraint",
				}
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(
						chat.Name, chat.Title, chat.Description,
						chat.OwnerID, chat.OwnerID, chat.UserIDUpdater, chat.ChatID,
					).
					WillReturnError(pgErr)
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeUniqueViolation,
		},
		{
			name: "Error - Bind Named Params",
			input: &model.UpdateChat{
				ChatID:        10,
				Name:          ptr("Updated Name"),
				Title:         ptr("Updated Title"),
				Description:   ptr("Updated Description"),
				OwnerID:       ptr(int64(2)),
				UserIDUpdater: 1,
			},
			mock: func(mock sqlmock.Sqlmock, chat *model.UpdateChat) {
				mock.ExpectQuery(`UPDATE chats SET *`).
					WithArgs(
						chat.Name, chat.Title, chat.Description,
						chat.OwnerID, chat.OwnerID, chat.UserIDUpdater, chat.ChatID,
					).
					WillReturnError(errDBBindError)
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

			err := repo.UpdateChat(context.Background(), tt.input)

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
