package pg

import (
	"chat-gateway/internal/entity"
	"context"
	"fmt"
	"strings"
)

// CreateChat accept user IDs and create direct chat.
func (p *PostgresRepository) CreateDirectChat(ctx context.Context, chat *entity.DirectChat) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	chatQuery := `
		INSERT INTO chats (type)
		VALUES ('direct')
		RETURNING id, created_at`

	err = tx.QueryRowxContext(ctx, chatQuery).Scan(&chat.ChatID, &chat.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert chat: %w", err)
	}

	membersQuery := `
		INSERT INTO chat_members (chat_id, user_id)
		VALUES ($1, $2), ($1, $3)`

	_, err = tx.ExecContext(ctx, membersQuery, chat.ChatID, chat.SenderID, chat.RecipientID)
	if err != nil {
		return fmt.Errorf("failed to insert chat members: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// CreateGroupChat accept user IDs, chat name, chat owner user_id and create group chat between several users.
func (p *PostgresRepository) CreateGroupChat(ctx context.Context, chat *entity.GroupChat) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
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
		return fmt.Errorf("failed to insert chat: %w", err)
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
		return fmt.Errorf("failed to insert chat members: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
