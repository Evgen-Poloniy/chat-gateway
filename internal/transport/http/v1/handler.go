package v1

import (
	"chat-gateway/internal/service"
)

type Handler struct {
	messenger service.MessengerService
}

func NewHandler(messenger service.MessengerService) *Handler {
	return &Handler{
		messenger: messenger,
	}
}
