package dto

import (
	"time"
)

// RegisterUser represents DTO for register new user data into the messenger database.
type RegisterUserReq struct {
	Username  string     `json:"username"`
	Email     *string    `json:"email,omitempty"`
	FirstName *string    `json:"first_name,omitempty"`
	LastName  *string    `json:"last_name,omitempty"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	Gender    *string    `json:"gender,omitempty"`
}

// UpdateUserReq represents DTO for update user data into the messenger database.
type UpdateUserReq struct {
	Username  *string    `json:"username,omitempty"`
	Email     *string    `json:"email,omitempty"`
	FirstName *string    `json:"first_name,omitempty"`
	LastName  *string    `json:"last_name,omitempty"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	Gender    *string    `json:"gender,omitempty"`
}

// UserInfo represents DTO response with data about user the messenger database.
type UserDataResp struct {
	UserID    int64      `json:"user_id"`
	Username  string     `json:"username"`
	Email     *string    `json:"email,omitempty"`
	FirstName *string    `json:"first_name,omitempty"`
	LastName  *string    `json:"last_name,omitempty"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	Gender    *string    `json:"gender,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
