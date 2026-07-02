package dto

// Error represents an error message
// swagger:model Error
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ResponseError wraps Error in "error" field
// swagger:model ResponseError
type ResponseError struct {
	Error Error `json:"error"`
}
