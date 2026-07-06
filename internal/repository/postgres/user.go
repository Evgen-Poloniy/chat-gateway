package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
	errs "github.com/Evgen-Poloniy/chat-gateway/pkg/errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// GetUserDataByUsername gets all data about user from the messenger database.
func (p *PostgresRepository) GetUserDataByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `
        SELECT id, username, email, first_name, last_name, birth_date, created_at, gender
        FROM users
        WHERE username = $1`

	var user entity.User

	if err := p.db.GetContext(ctx, &user, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with username '%s' not found: %w", username, errs.ErrRecordNotFound)
		}
		return nil, err
	}

	return &user, nil
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
		err = rows.Scan(&user.UserID, &user.CreatedAt)
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
