package repository

import (
	"chat-gateway/internal/entity"
	pg "chat-gateway/internal/repository/postgres"
	"context"

	"github.com/jmoiron/sqlx"
)

// MessengerRepository represents interface for work with the messenger database.
type MessengerRepository interface {
	// GetUserIdByUsername allows get username from messenger database and get unique used_id.
	GetUserIdByUsername(ctx context.Context, username string) (int64, error)

	// CreateUser allows create user into messenger database.
	CreateUser(ctx context.Context, user *entity.User) error

	// CreateChat accept user IDs and create direct chat.
	CreateDirectChat(ctx context.Context, chat *entity.DirectChat) error
}

type Repository struct {
	MessengerRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		MessengerRepository: pg.NewPostgresRepository(db),
	}
}
