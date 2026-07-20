package postgres_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Evgen-Poloniy/chat-gateway/internal/errs"
	"github.com/Evgen-Poloniy/chat-gateway/internal/model"
	pg "github.com/Evgen-Poloniy/chat-gateway/internal/repository/postgres"

	"github.com/google/uuid"
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
	chatID := uuid.New()
	userID1 := uuid.New()
	userID2 := uuid.New()

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
				ChatID:         chatID,
				ParticipantIDs: []uuid.UUID{userID1, userID2},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (id, type) VALUES ($1, 'direct') RETURNING created_at`)).
					WithArgs(chatID).
					WillReturnRows(sqlmock.NewRows([]string{"created_at"}).AddRow(now))

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO chat_members (chat_id, user_id) VALUES ($1, $2), ($1, $3)`)).
					WithArgs(chatID, userID1, userID2).
					WillReturnResult(sqlmock.NewResult(1, 2))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Error - Transaction Begin Failed",
			input: &model.DirectChat{
				ChatID:         chatID,
				ParticipantIDs: []uuid.UUID{userID1, userID2},
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
				ChatID:         chatID,
				ParticipantIDs: []uuid.UUID{userID1, userID2},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				pgErr := &pgconn.PgError{
					Code:    "23505",
					Message: "duplicate key value violates unique constraint",
					Detail:  "Key already exists.",
				}

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (id, type) VALUES ($1, 'direct') RETURNING created_at`)).
					WithArgs(chatID).
					WillReturnError(pgErr)

				mock.ExpectRollback()
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeUniqueViolation,
		},
		{
			name: "Error - Query Failed",
			input: &model.DirectChat{
				ChatID:         chatID,
				ParticipantIDs: []uuid.UUID{userID1, userID2},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (id, type) VALUES ($1, 'direct') RETURNING created_at`)).
					WithArgs(chatID).
					WillReturnError(errSyntaxError)

				mock.ExpectRollback()
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - Insert Members Failed",
			input: &model.DirectChat{
				ChatID:         chatID,
				ParticipantIDs: []uuid.UUID{userID1, userID2},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (id, type) VALUES ($1, 'direct') RETURNING created_at`)).
					WithArgs(chatID).
					WillReturnRows(sqlmock.NewRows([]string{"created_at"}).AddRow(now))

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO chat_members (chat_id, user_id) VALUES ($1, $2), ($1, $3)`)).
					WithArgs(chatID, userID1, userID2).
					WillReturnError(errFKViolation)

				mock.ExpectRollback()
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - Commit Failed",
			input: &model.DirectChat{
				ChatID:         chatID,
				ParticipantIDs: []uuid.UUID{userID1, userID2},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (id, type) VALUES ($1, 'direct') RETURNING created_at`)).
					WithArgs(chatID).
					WillReturnRows(sqlmock.NewRows([]string{"created_at"}).AddRow(now))

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO chat_members (chat_id, user_id) VALUES ($1, $2), ($1, $3)`)).
					WithArgs(chatID, userID1, userID2).
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
				assert.Equal(t, chatID, tt.input.ChatID)
				assert.Equal(t, now, tt.input.CreatedAt)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresRepository_CreateGroupChat(t *testing.T) {
	now := time.Now()
	chatID := uuid.New()
	ownerID := uuid.New()
	userID2 := uuid.New()
	userID3 := uuid.New()

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
				ChatID:         chatID,
				Name:           "Group 1",
				Title:          ptr("Group chat"),
				Description:    ptr("Group chat"),
				OwnerID:        ownerID,
				ParticipantIDs: []uuid.UUID{ownerID, userID2, userID3},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (id, type, name, title, description, owner_id) VALUES ($1, 'group', $2, $3, $4, $5) RETURNING created_at`)).
					WithArgs(chatID, "Group 1", "Group chat", "Group chat", ownerID).
					WillReturnRows(sqlmock.NewRows([]string{"created_at"}).AddRow(now))

				expectedQuery := regexp.QuoteMeta(`INSERT INTO chat_members (chat_id, user_id) VALUES ($1, $2), ($1, $3), ($1, $4)`)
				mock.ExpectExec(expectedQuery).
					WithArgs(chatID, ownerID, userID2, userID3).
					WillReturnResult(sqlmock.NewResult(1, 3))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Error - Transaction Begin Failed",
			input: &model.GroupChat{
				ChatID: chatID,
				Name:   "Group 2",
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
				ChatID:         chatID,
				Name:           "Duplicate",
				Title:          nil,
				Description:    nil,
				OwnerID:        ownerID,
				ParticipantIDs: []uuid.UUID{ownerID},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				pgErr := &pgconn.PgError{
					Code:    "23505",
					Message: "unique violation",
					Detail:  "detail",
				}

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (id, type, name, title, description, owner_id) VALUES ($1, 'group', $2, $3, $4, $5) RETURNING created_at`)).
					WithArgs(chatID, "Duplicate", nil, nil, ownerID).
					WillReturnError(pgErr)

				mock.ExpectRollback()
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeUniqueViolation,
		},
		{
			name: "Error - Foreign Key Violation",
			input: &model.GroupChat{
				ChatID:         chatID,
				Name:           "New Group",
				OwnerID:        ownerID,
				Title:          nil,
				Description:    nil,
				ParticipantIDs: []uuid.UUID{ownerID},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				pgErr := &pgconn.PgError{Code: "23503", Message: "fk violation", Detail: "owner_id does not exist"}
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (id, type, name, title, description, owner_id) VALUES ($1, 'group', $2, $3, $4, $5) RETURNING created_at`)).
					WithArgs(chatID, "New Group", nil, nil, ownerID).
					WillReturnError(pgErr)

				mock.ExpectRollback()
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeForeignKeyViolation,
		},
		{
			name: "Error - Insert Members Failed",
			input: &model.GroupChat{
				ChatID:         chatID,
				Name:           "Group",
				OwnerID:        ownerID,
				Title:          nil,
				Description:    nil,
				ParticipantIDs: []uuid.UUID{ownerID, userID2},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (id, type, name, title, description, owner_id) VALUES ($1, 'group', $2, $3, $4, $5) RETURNING created_at`)).
					WithArgs(chatID, "Group", nil, nil, ownerID).
					WillReturnRows(sqlmock.NewRows([]string{"created_at"}).AddRow(now))

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO chat_members (chat_id, user_id) VALUES ($1, $2), ($1, $3)`)).
					WithArgs(chatID, ownerID, userID2).
					WillReturnError(errMemberInsertError)

				mock.ExpectRollback()
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - Query Failed",
			input: &model.GroupChat{
				ChatID:         chatID,
				Name:           "New Group",
				OwnerID:        ownerID,
				Title:          nil,
				Description:    nil,
				ParticipantIDs: []uuid.UUID{ownerID},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (id, type, name, title, description, owner_id) VALUES ($1, 'group', $2, $3, $4, $5) RETURNING created_at`)).
					WithArgs(chatID, "New Group", nil, nil, ownerID).
					WillReturnError(errDBQueryFailed)

				mock.ExpectRollback()
			},
			wantErr:           true,
			expectedErrorCode: errs.CodeQueryError,
		},
		{
			name: "Error - Commit Failed",
			input: &model.GroupChat{
				ChatID:         chatID,
				Name:           "Group",
				OwnerID:        ownerID,
				Title:          nil,
				Description:    nil,
				ParticipantIDs: []uuid.UUID{ownerID},
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO chats (id, type, name, title, description, owner_id) VALUES ($1, 'group', $2, $3, $4, $5) RETURNING created_at`)).
					WithArgs(chatID, "Group", nil, nil, ownerID).
					WillReturnRows(sqlmock.NewRows([]string{"created_at"}).AddRow(now))

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO chat_members (chat_id, user_id) VALUES ($1, $2)`)).
					WithArgs(chatID, ownerID).
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
				assert.Equal(t, chatID, tt.input.ChatID)
				assert.Equal(t, now, tt.input.CreatedAt)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresRepository_GetChatsByUserID(t *testing.T) {
	now := time.Now()
	userID := uuid.New()
	chatID := uuid.New()
	ownerID := uuid.New()

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
		userID            uuid.UUID
		limit             int
		offset            int
		mock              func(mock sqlmock.Sqlmock)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name:   "Success",
			userID: userID,
			limit:  10,
			offset: 0,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "type", "name", "title", "description", "created_at", "owner_id"}).
					AddRow(chatID, "group", "group", nil, nil, now, ownerID)

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(userID, 10, 0).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:   "Error - Query Failed",
			userID: userID,
			limit:  10,
			offset: 0,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(userID, 10, 0).
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

				assert.Equal(t, chatID, chats[0].ChatID)
				assert.Equal(t, "group", chats[0].ChatType)
				assert.Equal(t, "group", *chats[0].Name)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresRepository_UpdateGroupChat(t *testing.T) {
	now := time.Now()
	chatID := uuid.New()
	updaterID := uuid.New()
	newOwnerID := uuid.New()

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
		input             *model.UpdateGroupChat
		mock              func(mock sqlmock.Sqlmock, chat *model.UpdateGroupChat)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name: "Success",
			input: &model.UpdateGroupChat{
				ChatID:        chatID,
				Name:          ptr("Updated Name"),
				Title:         ptr("Updated Title"),
				Description:   ptr("Updated Description"),
				OwnerID:       ptr(newOwnerID),
				UserIDUpdater: updaterID,
			},
			mock: func(mock sqlmock.Sqlmock, chat *model.UpdateGroupChat) {
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
			input: &model.UpdateGroupChat{
				ChatID:        chatID,
				Name:          ptr("Updated Name"),
				Title:         ptr("Updated Title"),
				Description:   ptr("Updated Description"),
				OwnerID:       ptr(newOwnerID),
				UserIDUpdater: updaterID,
			},
			mock: func(mock sqlmock.Sqlmock, chat *model.UpdateGroupChat) {
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
			input: &model.UpdateGroupChat{
				ChatID:        chatID,
				Name:          ptr("Updated Name"),
				Title:         ptr("Updated Title"),
				Description:   ptr("Updated Description"),
				OwnerID:       ptr(newOwnerID),
				UserIDUpdater: updaterID,
			},
			mock: func(mock sqlmock.Sqlmock, chat *model.UpdateGroupChat) {
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
			input: &model.UpdateGroupChat{
				ChatID:        chatID,
				Name:          ptr("Updated Name"),
				Title:         ptr("Updated Title"),
				Description:   ptr("Updated Description"),
				OwnerID:       ptr(newOwnerID),
				UserIDUpdater: updaterID,
			},
			mock: func(mock sqlmock.Sqlmock, chat *model.UpdateGroupChat) {
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

			err := repo.UpdateGroupChat(context.Background(), tt.input)

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

func TestPostgresRepository_IsUserInChat(t *testing.T) {
	chatID := uuid.New()
	userID := uuid.New()

	query := `
    SELECT EXISTS (
        SELECT 1
        FROM chat_members
        WHERE chat_id = $1 AND user_id = $2
    )`

	tests := []struct {
		name              string
		chatID            uuid.UUID
		userID            uuid.UUID
		mock              func(mock sqlmock.Sqlmock)
		wantErr           bool
		expectedErrorCode errs.ErrCode
	}{
		{
			name:   "Success - User Is Member",
			chatID: chatID,
			userID: userID,
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"exists"}).
					AddRow(true)

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(chatID, userID).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name:   "Error - Query Failed",
			chatID: chatID,
			userID: userID,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(chatID, userID).
					WillReturnError(errSyntaxError)
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

			isMember, err := repo.IsUserInChat(context.Background(), tt.chatID, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				var appErr *errs.AppError
				if assert.ErrorAs(t, err, &appErr) {
					assert.Equal(t, tt.expectedErrorCode, appErr.Code)
				}
				assert.False(t, isMember)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
