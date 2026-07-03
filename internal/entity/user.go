package entity

import "time"

// User represents user info messenger
type User struct {
	ID        int64      `db:"id"`
	Username  string     `db:"username"`
	Email     *string    `db:"email"`
	FirstName *string    `db:"first_name"`
	LastName  *string    `db:"last_name"`
	BirthDate *time.Time `db:"birth_date"`
	CreatedAt *time.Time `db:"created_at"`
}
