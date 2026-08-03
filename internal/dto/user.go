package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// CasdoorWebhookReq is the top-level payload structure sent by Casdoor webhook.
type CasdoorWebhookReq struct {
	Action       string          `json:"action"`
	ExtendedUser json.RawMessage `json:"extendedUser"`
}

// SignUp contains user profile data for "signup" event.
type SignUpUser struct {
	UserID    string  `json:"id"`
	Username  string  `json:"name"`
	Email     *string `json:"email,omitempty"`
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
	Birthday  *string `json:"birthday,omitempty"`
	Gender    *string `json:"gender,omitempty"`
}

// Update contains user profile data for "update" event.
type UpdateUser struct {
	UserID    string  `json:"id"`
	Username  *string `json:"name,omitempty"`
	Email     *string `json:"email,omitempty"`
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
	Birthday  *string `json:"birthday,omitempty"`
	Gender    *string `json:"gender,omitempty"`
}

// Delete contains minimal user identification data for "delete-user" event.
type DeleteUser struct {
	UserID string `json:"id"`
}

// UserInfo represents DTO response with data about user the messenger database.
type UserDataResp struct {
	UserID    uuid.UUID  `json:"user_id"`
	Username  string     `json:"username"`
	Email     *string    `json:"email,omitempty"`
	FirstName *string    `json:"first_name,omitempty"`
	LastName  *string    `json:"last_name,omitempty"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	Gender    *string    `json:"gender,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
