package dto

import (
	"time"
)

// RegisterUser represents DTO for register new user info messenger.
type RegisterUser struct {
	Username  string     `json:"username"`
	Email     *string    `json:"email,omitempty"`
	FirstName *string    `json:"first_name,omitempty"`
	LastName  *string    `json:"last_name,omitempty"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	Gender    *string    `json:"gender,omitempty"`
}

// UserInfo represents request with data about user.
type UserData struct {
	UserID    int64      `json:"user_id"`
	Username  string     `json:"username,omitempty"`
	Email     *string    `json:"email,omitempty"`
	FirstName *string    `json:"first_name,omitempty"`
	LastName  *string    `json:"last_name,omitempty"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	Gender    *string    `json:"gender,omitempty"`
}
