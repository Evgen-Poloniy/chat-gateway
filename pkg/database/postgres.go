package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Evgen-Poloniy/chat-gateway/internal/config"

	"github.com/jmoiron/sqlx"
)

// NewPostgreSQL opens connection with the database.
func NewPostgreSQL(config *config.PostgresConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode,
	)

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("error when open connection with the database: %w", err)
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("the database is not available: %w", err)
	}

	return db, nil
}
