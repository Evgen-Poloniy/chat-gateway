package pg

import "github.com/jmoiron/sqlx"

// PostgresRepository represents implementation of interface MessengerRepository
type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}
