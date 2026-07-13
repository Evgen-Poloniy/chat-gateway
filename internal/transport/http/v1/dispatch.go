package v1

import (
	"net/http"

	"github.com/Evgen-Poloniy/chat-gateway/internal/dto"
	errs "github.com/Evgen-Poloniy/chat-gateway/pkg/errors"
	"github.com/gin-gonic/gin"
)

// DispatchMessage dispatches message for clients
func (h *Handler) DispatchMessage(c *gin.Context) {
	var req dto.DispatchMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(&errs.HttpError{
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
			Message:    err.Error(),
		})
		return
	}

	offlineIDs := make([]int64, 0, len(req.RecipientIDs))

	message := dto.DispatchMessage{
		SenderID: req.SenderID,
		Message:  req.Message,
	}

	for _, userID := range req.RecipientIDs {
		if err := h.wsHub.SendToUser(userID, message); err != nil {
			offlineIDs = append(offlineIDs, userID)
		}
	}

	resp := dto.DispatchMessageResp{
		OfflineUserIDs: offlineIDs,
	}

	c.JSON(http.StatusOK, dto.DataResp{Data: resp})
}
