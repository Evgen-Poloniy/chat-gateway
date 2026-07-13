package v1

import (
	"github.com/Evgen-Poloniy/chat-gateway/internal/transport/http/ws"
	dispatchv1 "github.com/Evgen-Poloniy/messenger-contracts/gen/go/dispatch/v1"
)

type Handler struct {
	dispatchv1.UnimplementedDispatchServiceServer
	wsHub *ws.Hub
}

func NewHandler(wsHub *ws.Hub) *Handler {
	return &Handler{
		wsHub: wsHub,
	}
}
