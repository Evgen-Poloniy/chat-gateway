package pg

import (
	"chat-gateway/internal/entity"
	"context"
	"fmt"
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
