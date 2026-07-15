package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Evgen-Poloniy/chat-gateway/internal/model"
	errs "github.com/Evgen-Poloniy/chat-gateway/pkg/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

// CreateDirectChat accept user IDs and create direct chat.
func (p *PostgresRepository) CreateDirectChat(ctx context.Context, chat *model.DirectChat) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return errs.NewAppError(
			errs.CodeTransactionError,
			"database error: failed to begin transaction",
			fmt.Errorf("database error: %v", err),
		)
	}
	defer tx.Rollback()

	chatQuery := `
		INSERT INTO chats (id, type)
		VALUES ($1, 'direct')
		RETURNING created_at`

	if err = tx.QueryRowxContext(ctx, chatQuery, chat.ChatID).Scan(&chat.CreatedAt); err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23505" {
				return errs.NewAppError(
					errs.CodeUniqueViolation,
					"database error: "+pgErr.Message,
					fmt.Errorf("database error: %s", pgErr.Detail),
				)
			}
		}

		return errs.NewAppError(
			errs.CodeQueryError,
			"database error: failed to insert values into table",
			fmt.Errorf("database error: %v", err),
		)
	}

	membersQuery := `
		INSERT INTO chat_members (chat_id, user_id)
		VALUES ($1, $2), ($1, $3)`

	_, err = tx.ExecContext(ctx, membersQuery, chat.ChatID, chat.ParticipantIDs[0], chat.ParticipantIDs[1])
	if err != nil {
		return errs.NewAppError(
			errs.CodeQueryError,
			"database error: failed to insert values into table",
			fmt.Errorf("database error: %v", err),
		)
	}

	if err := tx.Commit(); err != nil {
		return errs.NewAppError(
			errs.CodeTransactionError,
			"database error: failed to commit transaction",
			fmt.Errorf("database error: %v", err),
		)
	}

	return nil
}

// CreateGroupChat accept user IDs, chat name, chat owner user_id and create group chat between several users.
func (p *PostgresRepository) CreateGroupChat(ctx context.Context, chat *model.GroupChat) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return errs.NewAppError(
			errs.CodeTransactionError,
			"database error: failed to begin transaction",
			fmt.Errorf("database error: %v", err),
		)
	}
	defer tx.Rollback()

	chatQuery := `
		INSERT INTO chats (id, type, name, title, description, owner_id)
		VALUES ($1, 'group', $2, $3, $4, $5)
		RETURNING created_at`

	err = tx.QueryRowxContext(
		ctx, chatQuery, chat.ChatID, chat.Name, chat.Title, chat.Description, chat.OwnerID,
	).Scan(&chat.CreatedAt)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			switch pgErr.Code {
			case "23505":
				return errs.NewAppError(
					errs.CodeUniqueViolation,
					"database error: "+pgErr.Message,
					fmt.Errorf("database error: %s", pgErr.Detail),
				)

			case "23503":
				return errs.NewAppError(
					errs.CodeForeignKeyViolation,
					"database error: owner_id does not exist",
					fmt.Errorf("database error: %s", pgErr.Detail),
				)
			}
		}

		return errs.NewAppError(
			errs.CodeQueryError,
			"database error: failed to insert values into table",
			fmt.Errorf("database error: %v", err),
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
			errs.CodeQueryError,
			"database error: failed to insert values into table",
			fmt.Errorf("database error: %v", err),
		)
	}

	if err := tx.Commit(); err != nil {
		return errs.NewAppError(
			errs.CodeTransactionError,
			"database error: failed to commit transaction",
			fmt.Errorf("database error: %v", err),
		)
	}

	return nil
}

// GetChatsByUserID gets chat by user_id with limits and offset
func (p *PostgresRepository) GetChatsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Chat, error) {
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
		LIMIT $2 OFFSET $3
	`

	var chats []model.Chat

	err := p.db.SelectContext(ctx, &chats, query, userID, limit, offset)
	if err != nil {
		return nil, errs.NewAppError(
			errs.CodeQueryError,
			fmt.Sprintf("database error: failed to get chats by user id: %s", userID),
			fmt.Errorf("database error: %v", err),
		)
	}

	if len(chats) == 0 {
		return nil, errs.NewAppError(
			errs.CodeChatNotFound,
			fmt.Sprintf("database error: records about user chats with user_id '%s' not found", userID),
			fmt.Errorf("database error: records about user chats with user_id '%s' not found", userID),
		)
	}

	return chats, nil
}

// UpdateGroupChat updates data about chat like name, title, description, owner
func (p *PostgresRepository) UpdateGroupChat(ctx context.Context, chat *model.UpdateGroupChat) error {
	query := `
		UPDATE chats
		SET
			name = COALESCE(:name, name),
			title = COALESCE(:title, title),
			description = COALESCE(:description, description),
			owner_id = CASE
				WHEN :owner_id IS NULL THEN owner_id
				WHEN :owner_id = :user_id_updater THEN owner_id
				ELSE owner_id
			END
		WHERE id = :id
		RETURNING name, title, description, created_at, owner_id
	`

	boundQuery, args, err := p.db.BindNamed(query, chat)
	if err != nil {
		return errs.NewAppError(
			errs.CodeQueryError,
			"database error: failed to bind named params",
			fmt.Errorf("database error: %v", err),
		)
	}

	if err := p.db.QueryRowxContext(ctx, boundQuery, args...).StructScan(chat); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.NewAppError(
				errs.CodeChatNotFound,
				fmt.Sprintf("database error: record about chat with chat_id '%s' not found", chat.ChatID),
				fmt.Errorf("database error: record about chat with chat_id '%s' not found", chat.ChatID),
			)
		}

		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23505" {
				return errs.NewAppError(
					errs.CodeUniqueViolation,
					"database error: "+pgErr.Message,
					fmt.Errorf("database error: %s", pgErr.Detail),
				)
			}
		}

		return errs.NewAppError(
			errs.CodeQueryError,
			"database error: failed to insert values into table",
			fmt.Errorf("database error: %v", err),
		)
	}

	return nil
}
