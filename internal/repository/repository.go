package repository

import (
	pg "chat-gateway/internal/repository/postgres"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// MessengerRepository represents interface for work with the messenger database
type MessengerRepository interface {
}

type Repository struct {
	MessengerRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		MessengerRepository: pg.NewPostgresRepository(db),
	}
}
