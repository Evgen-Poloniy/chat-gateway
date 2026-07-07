package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
	errs "github.com/Evgen-Poloniy/chat-gateway/pkg/errors"
	"github.com/jackc/pgx/v5/pgconn"

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
			return nil, errs.NewAppError(
				"NOT_FOUND",
				fmt.Sprintf("database error: record with username '%s' not found", username),
				fmt.Sprintf("database error: record with username '%s' not found", username),
				errs.ErrRecordNotFound,
			)
		}

		return nil, errs.NewAppError(
			"QUERY_ERROR",
			"database error: query error",
			fmt.Sprintf("database error: %v", err),
			errs.ErrQuery,
		)
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
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&user.UserID, &user.CreatedAt)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errs.ErrRecordNotFound
			}

			return errs.NewAppError(
				"QUERY_ERROR",
				"database error: failed to scan rows",
				fmt.Sprintf("database error: %v", err),
				errs.ErrQuery,
			)
		}

		return nil
	}

	if err = rows.Err(); err != nil {
		return errs.NewAppError(
			"QUERY_ERROR",
			"database error: rows error",
			fmt.Sprintf("database error: %v", err),
			errs.ErrQuery,
		)
	}

	return errs.NewAppError(
		"NOT_FOUND",
		"database error: record not found",
		"database error: record not found",
		errs.ErrRecordNotFound,
	)
}

// UpdateUser updates data about user into messenger database.
func (p *PostgresRepository) UpdateUser(ctx context.Context, user *entity.UpdateUser) error {
	query := `
        UPDATE users
        SET
            username = COALESCE(:username, username),
            email = COALESCE(:email, email),
            first_name = COALESCE(:first_name, first_name),
            last_name = COALESCE(:last_name, last_name),
            birth_date = COALESCE(:birth_date, birth_date),
            gender = COALESCE(:gender, gender)
        WHERE id = :id
    `

	if _, err := p.db.NamedExecContext(ctx, query, user); err != nil {
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

	return nil
}
