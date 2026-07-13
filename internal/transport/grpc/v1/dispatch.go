package v1

import (
	"context"

	"github.com/Evgen-Poloniy/chat-gateway/internal/dto"
	dispatchv1 "github.com/Evgen-Poloniy/messenger-contracts/gen/go/dispatch/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) DispatchMessage(ctx context.Context, req *dispatchv1.DispatchMessageRequest) (*dispatchv1.DispatchMessageResponse, error) {
	if req.GetSenderId() <= 0 || len(req.GetRecipientIds()) == 0 || req.GetMessage() == "" {
		return nil, status.Error(codes.InvalidArgument, "bad request: missing required fields")
	}

	message := dto.DispatchMessage{
		SenderID: req.GetSenderId(),
		Message:  req.GetMessage(),
	}

	offlineIDs := make([]int64, 0, len(req.GetRecipientIds()))

	for _, userID := range req.GetRecipientIds() {
		if err := h.wsHub.SendToUser(userID, message); err != nil {
			offlineIDs = append(offlineIDs, userID)
		}
	}

	resp := &dispatchv1.DispatchMessageResponse{
		OfflineUserIds: offlineIDs,
	}

	return resp, nil
}
