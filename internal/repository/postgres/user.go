package pg

import (
	"chat-gateway/internal/entity"
	errs "chat-gateway/pkg/errors"
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// GetUserIdByUsername allows get username from messenger database and get unique used_id.
func (p *PostgresRepository) GetUserIdByUsername(ctx context.Context, username string) (int64, error) {
	query := "SELECT id FROM users WHERE username LIKE '$1'"

	var user_id int64

	if err := p.db.GetContext(ctx, &user_id, query, username); err != nil {
		return 0, errs.ErrRecordNotFound
	}

	return user_id, nil
}

// CreateUser allows create user into messenger database.
func (p *PostgresRepository) CreateUser(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (username, email, first_name, last_name, birth_date, gender)
		VALUES (:username, :email, :first_name, :last_name, :birth_date, :gender)
		RETURNING id, created_at
	`

	rows, err := p.db.NamedQueryContext(ctx, query, user)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&user.ID, &user.CreatedAt)
		if err != nil {
			return err
		}

		return nil
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return errs.ErrRecordNotFound
}
