package dto

// Error represents an error message
// swagger:model Error
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse wraps Error in "error" field
// swagger:model ErrorResponse
type ErrorResponse struct {
	Error Error `json:"error"`
}
