package pg

import (
	"chat-gateway/internal/entity"
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// CreateUser allows create user into messenger database and get unique used_id
func (p *PostgresRepository) CreateUser(ctx context.Context, user *entity.User) (int64, error) {
	query := `
		INSERT INTO users (username, email, first_name, last_name, birth_date)
		VALUES (:username, :email, :first_name, :last_name, :birth_date)
		RETURNING id
	`

	var user_id int64

	rows, err := p.db.NamedQueryContext(ctx, query, user)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(user)
		if err != nil {
			return 0, err
		}
	}

	return user_id, nil
}
