package model

import "time"

// User represents user data into messenger.
type User struct {
	UserID    int64      `db:"id"`
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
	UserID    int64      `db:"id"`
	Username  *string    `db:"username"`
	Email     *string    `db:"email"`
	FirstName *string    `db:"first_name"`
	LastName  *string    `db:"last_name"`
	BirthDate *time.Time `db:"birth_date"`
	Gender    *string    `db:"gender"`
	CreatedAt time.Time  `db:"created_at"`
}
