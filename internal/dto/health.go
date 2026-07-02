package dto

// Health represents API health status
// swagger:model Health
type Health struct {
	Status string `json:"status"`
}
