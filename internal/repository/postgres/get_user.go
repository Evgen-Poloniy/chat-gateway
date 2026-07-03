package pg

import (
	errs "chat-gateway/pkg/errors"
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// GetUserIdByUsername allows get username from messenger database and get unique used_id
func (p *PostgresRepository) GetUserIdByUsername(ctx context.Context, username string) (int64, error) {
	query := "SELECT id FROM users WHERE username LIKE '$1'"

	var user_id int64

	if err := p.db.GetContext(ctx, &user_id, query, username); err != nil {
		return 0, errs.ErrRecordNotFound
	}

	return user_id, nil
}
