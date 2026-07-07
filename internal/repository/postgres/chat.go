package pg

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
	errs "github.com/Evgen-Poloniy/chat-gateway/pkg/errors"
	"github.com/jackc/pgx/v5/pgconn"
)

// CreateChat accept user IDs and create direct chat.
func (p *PostgresRepository) CreateDirectChatReq(ctx context.Context, chat *entity.DirectChat) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return errs.NewAppError(
			"TRANSACTION_ERROR",
			"database error: failed to begin transaction",
			fmt.Sprintf("database error: %v", err),
			errs.ErrFailedToBeginTransaction,
		)
	}
	defer tx.Rollback()

	chatQuery := `
		INSERT INTO chats (type)
		VALUES ('direct')
		RETURNING id, created_at`

	err = tx.QueryRowxContext(ctx, chatQuery).Scan(&chat.ChatID, &chat.CreatedAt)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23505" {
				return errs.NewAppError(
					"UNIQUE_VIOLATION",
					"database error: "+pgErr.Message,
					fmt.Sprintf("database error: %s", pgErr.Detail),
					errs.ErrUniqueViolation,
				)
			}
		}

		return errs.NewAppError(
			"QUERY_ERROR",
			"database error: failed to insert values into table",
			fmt.Sprintf("database error: %v", err),
			errs.ErrQuery,
		)
	}

	membersQuery := `
		INSERT INTO chat_members (chat_id, user_id)
		VALUES ($1, $2), ($1, $3)`

	_, err = tx.ExecContext(ctx, membersQuery, chat.ChatID, chat.SenderID, chat.RecipientID)
	if err != nil {
		return errs.NewAppError(
			"QUERY_ERROR",
			"database error: failed to insert values into table",
			fmt.Sprintf("database error: %v", err),
			errs.ErrQuery,
		)
	}

	if err := tx.Commit(); err != nil {
		return errs.NewAppError(
			"TRANSACTION_ERROR",
			"database error: failed to commit transaction",
			fmt.Sprintf("database error: %v", err),
			errs.ErrFailedToCommitTransaction,
		)
	}

	return nil
}

// CreateGroupChatReq accept user IDs, chat name, chat owner user_id and create group chat between several users.
func (p *PostgresRepository) CreateGroupChatReq(ctx context.Context, chat *entity.GroupChat) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return errs.NewAppError(
			"TRANSACTION_ERROR",
			"database error: failed to begin transaction",
			fmt.Sprintf("database error: %v", err),
			errs.ErrFailedToBeginTransaction,
		)
	}
	defer tx.Rollback()

	chatQuery := `
		INSERT INTO chats (type, name, title, description, owner_id)
		VALUES ('group', $1, $2, $3, $4)
		RETURNING id, created_at`

	err = tx.QueryRowxContext(
		ctx, chatQuery, chat.Name, chat.Title, chat.Description, chat.OwnerID,
	).Scan(&chat.ChatID, &chat.CreatedAt)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23505" {
				return errs.NewAppError(
					"UNIQUE_VIOLATION",
					"database error: "+pgErr.Message,
					fmt.Sprintf("database error: %s", pgErr.Detail),
					errs.ErrUniqueViolation,
				)
			}
		}

		return errs.NewAppError(
			"QUERY_ERROR",
			"database error: failed to insert values into table",
			fmt.Sprintf("database error: %v", err),
			errs.ErrQuery,
		)
	}

	values := make([]string, 0, len(chat.ParticipantIDs))
	args := make([]interface{}, 0, 1+len(chat.ParticipantIDs))

	args = append(args, chat.ChatID)

	for i, userID := range chat.ParticipantIDs {
		values = append(values, fmt.Sprintf("($1, $%d)", i+2))
		args = append(args, userID)
	}

	membersQuery := "INSERT INTO chat_members (chat_id, user_id) VALUES " + strings.Join(values, ", ")

	_, err = tx.ExecContext(ctx, membersQuery, args...)
	if err != nil {
		return errs.NewAppError(
			"QUERY_ERROR",
			"database error: failed to insert values into table",
			fmt.Sprintf("database error: %v", err),
			errs.ErrQuery,
		)
	}

	if err := tx.Commit(); err != nil {
		return errs.NewAppError(
			"TRANSACTION_ERROR",
			"database error: failed to commit transaction",
			fmt.Sprintf("database error: %v", err),
			errs.ErrFailedToCommitTransaction,
		)
	}

	return nil
}
