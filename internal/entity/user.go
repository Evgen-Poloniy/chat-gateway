package entity

import (
	"time"
)

// User represents user data into messenger.
type User struct {
	UserID    int64      `validate:"required,gt=0"`
	Username  string     `validate:"required,min=3,max=32,alphanum"`
	Email     *string    `validate:"omitempty,email,max=64"`
	FirstName *string    `validate:"omitempty,min=1,max=100"`
	LastName  *string    `validate:"omitempty,min=1,max=100"`
	BirthDate *time.Time `validate:"omitempty"`
	Gender    *string    `validate:"omitempty,oneof=man woman"`
	CreatedAt time.Time  `validate:"-"`
}

// UpdateUser represents user updating data into messenger.
type UpdateUser struct {
	UserID    int64      `validate:"required,gt=0"`
	Username  *string    `validate:"omitempty,min=3,max=32,alphanum"`
	Email     *string    `validate:"omitempty,email,max=64"`
	FirstName *string    `validate:"omitempty,min=1,max=100"`
	LastName  *string    `validate:"omitempty,min=1,max=100"`
	BirthDate *time.Time `validate:"omitempty"`
	Gender    *string    `validate:"omitempty,oneof=man woman"`
	CreatedAt time.Time  `validate:"-"`
}
