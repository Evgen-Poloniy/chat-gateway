package ws

import (
	"github.com/Evgen-Poloniy/chat-gateway/internal/errs"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var wsCloseCodeMap = map[errs.ErrCode]int{
	errs.CodeFailedEventChannel:   websocket.CloseInternalServerErr,
	errs.CodeDeserializationError: websocket.CloseInternalServerErr,
	errs.CodeDeliveryFailed:       websocket.CloseInternalServerErr,
	errs.CodeEmptyUserIDs:         websocket.ClosePolicyViolation,
	errs.CodeValidationError:      websocket.ClosePolicyViolation,
}

type Handler struct {
	wsHub  *Hub
	logger *logrus.Logger
}

func NewHandler(wsHub *Hub, logger *logrus.Logger) *Handler {
	return &Handler{
		wsHub:  wsHub,
		logger: logger,
	}
}
