package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents user data into messenger.
type User struct {
	UserID    uuid.UUID  `db:"id"`
	Username  string     `db:"username"`
	Email     *string    `db:"email"`
	FirstName *string    `db:"first_name"`
	LastName  *string    `db:"last_name"`
	BirthDate *time.Time `db:"birth_date"`
	Gender    *string    `db:"gender"`
	CreatedAt time.Time  `db:"created_at"`
}

// UpdateUser represents user updating data into messenger.
type UpdateUser struct {
	UserID    uuid.UUID  `db:"id"`
	Username  *string    `db:"username"`
	Email     *string    `db:"email"`
	FirstName *string    `db:"first_name"`
	LastName  *string    `db:"last_name"`
	BirthDate *time.Time `db:"birth_date"`
	Gender    *string    `db:"gender"`
	CreatedAt time.Time  `db:"created_at"`
}
