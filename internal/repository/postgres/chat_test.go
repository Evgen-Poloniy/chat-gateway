package pg_test

import (
	"context"
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

func TestPostgresRepository_CreateDirectChat(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		input   *entity.DirectChat
		mock    func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Success",
			input: &entity.DirectChat{
				SenderID:    1,
				RecipientID: 2,
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
			name: "Error - Unique Violation in chats",
			input: &entity.DirectChat{
				SenderID:    1,
				RecipientID: 2,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				pgErr := &pgconn.PgError{
					Code:    "23505",
					Message: "duplicate key value violates unique constraint",
					Detail:  "Key already exists.",
				}

				mock.ExpectQuery(`.*INSERT INTO chats.*`).WillReturnError(pgErr)

				mock.ExpectRollback()
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

			err := repo.CreateDirectChat(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
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
		name    string
		input   *entity.GroupChat
		mock    func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Success",
			input: &entity.GroupChat{
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
			input: &entity.GroupChat{
				Name: "Group 2",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(context.DeadlineExceeded)
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

			err := repo.CreateGroupChat(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
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
		name    string
		userID  int64
		limit   int
		offset  int
		mock    func(mock sqlmock.Sqlmock)
		wantErr bool
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
			name:   "Query Error",
			userID: 1,
			limit:  10,
			offset: 0,
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(int64(1), 10, 0).
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

			chats, err := repo.GetChatsByUserID(context.Background(), tt.userID, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
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
