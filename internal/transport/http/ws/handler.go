package ws

import (
	"errors"
	"fmt"

	errs "github.com/Evgen-Poloniy/chat-gateway/pkg/errors"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var wsCloseCodeMap = map[errs.ErrCode]int{
	errs.CodeFailedEventChannel:   websocket.CloseInternalServerErr,
	errs.CodeDeserializationError: websocket.CloseInternalServerErr,
	errs.ErrCodeDeliveryFailed:    websocket.CloseInternalServerErr,
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

func logInfo(logger *logrus.Logger, id string, serverAddress string, message string) {
	entry := logger.WithFields(map[string]interface{}{
		"id":      id,
		"address": serverAddress,
	})
	entry.Info(message)
}

func logWarn(logger *logrus.Logger, id string, serverAddress string, message string) {
	entry := logger.WithFields(map[string]interface{}{
		"id":      id,
		"address": _serverAddress,
	})
	entry.Warn(message)
}

func logError(logger *logrus.Logger, id string, serverAddress string, err error) {
	code := "unknown_internal_error"
	closeCode := websocket.CloseInternalServerErr

	if appErr, ok := errors.AsType[*errs.AppError](err); ok {
		if str, exists := errs.MapToString[appErr.Code]; exists {
			code = str
		} else {
			code = fmt.Sprintf("close_code_%d", appErr.Code)
		}

		if wsCloseCode, exists := wsCloseCodeMap[appErr.Code]; exists {
			closeCode = wsCloseCode
		}
	}

	logger.WithFields(map[string]interface{}{
		"id":         id,
		"address":    serverAddress,
		"code":       code,
		"close_code": closeCode,
	}).Error(err)
}
