package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Evgen-Poloniy/chat-gateway/internal/model"
	errs "github.com/Evgen-Poloniy/chat-gateway/pkg/errors"
	"github.com/jackc/pgx/v5/pgconn"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// GetUserDataByUsername gets all data about user from the messenger database.
func (p *PostgresRepository) GetUserDataByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
        SELECT id, username, email, first_name, last_name, birth_date, created_at, gender
        FROM users
        WHERE username = $1`

	var user model.User

	if err := p.db.GetContext(ctx, &user, query, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewAppError(
				errs.CodeUserNotFound,
				fmt.Sprintf("database error: record about user with username '%s' not found", username),
				fmt.Errorf("database error: record about user with username '%s' not found", username),
			)
		}

		return nil, errs.NewAppError(
			errs.CodeQueryError,
			"database error: query error",
			fmt.Errorf("database error: %v", err),
		)
	}

	return &user, nil
}

// GetUserDataByUserID gets all data about user from the messenger database by user_id.
func (p *PostgresRepository) GetUserDataByUserID(ctx context.Context, userID int64) (*model.User, error) {
	query := `
        SELECT id, username, email, first_name, last_name, birth_date, created_at, gender
        FROM users
        WHERE id = $1`

	var user model.User

	if err := p.db.GetContext(ctx, &user, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewAppError(
				errs.CodeUserNotFound,
				fmt.Sprintf("database error: record about user with user_id '%d' not found", userID),
				fmt.Errorf("database error: record about user with user_id '%d' not found", userID),
			)
		}

		return nil, errs.NewAppError(
			errs.CodeQueryError,
			"database error: query error",
			fmt.Errorf("database error: %v", err),
		)
	}

	return &user, nil
}

// CreateUser allows create user into messenger database.
func (p *PostgresRepository) CreateUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (username, email, first_name, last_name, birth_date, gender)
		VALUES (:username, :email, :first_name, :last_name, :birth_date, :gender)
		RETURNING id, created_at
	`

	boundQuery, args, err := p.db.BindNamed(query, user)
	if err != nil {
		return errs.NewAppError(
			errs.CodeQueryError,
			"database error: failed to bind named params",
			fmt.Errorf("database error: %v", err),
		)
	}

	if err = p.db.QueryRowxContext(ctx, boundQuery, args...).Scan(&user.UserID, &user.CreatedAt); err != nil {
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

// UpdateUser updates data about user into messenger database.
func (p *PostgresRepository) UpdateUser(ctx context.Context, user *model.UpdateUser) error {
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
		RETURNING created_at
    `

	boundQuery, args, err := p.db.BindNamed(query, user)
	if err != nil {
		return errs.NewAppError(
			errs.CodeQueryError,
			"database error: failed to bind named params",
			fmt.Errorf("database error: %v", err),
		)
	}

	if err := p.db.QueryRowxContext(ctx, boundQuery, args...).Scan(&user.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.NewAppError(
				errs.CodeUserNotFound,
				fmt.Sprintf("database error: record about user with user_id '%d' not found", user.UserID),
				fmt.Errorf("database error: record about user with user_id '%d' not found", user.UserID),
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
