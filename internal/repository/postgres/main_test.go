package pg_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var (
	errDBQueryFailed = errors.New("db query failed")
	errDBBindError   = errors.New("db bind error")
)

func setupMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	cleanup := func() {
		mockDB.Close()
	}

	return sqlxDB, mock, cleanup
}

func ptr[T any](v T) *T {
	return &v
}
