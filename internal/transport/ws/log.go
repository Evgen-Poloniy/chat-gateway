package ws

import (
	"errors"

	"github.com/Evgen-Poloniy/chat-gateway/internal/errs"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func logInfo(logger *logrus.Logger, serverAddress string, message string) {
	entry := logger.WithFields(map[string]interface{}{
		"address": serverAddress,
	})
	entry.Info(message)
}

func logWarn(logger *logrus.Logger, serverAddress string, message string) {
	entry := logger.WithFields(map[string]interface{}{
		"address": serverAddress,
	})
	entry.Warn(message)
}

func logError(logger *logrus.Logger, serverAddress string, err error) {
	code := "unknown_internal_error"
	closeCode := websocket.CloseInternalServerErr

	if appErr, ok := errors.AsType[*errs.AppError](err); ok {
		code = string(appErr.Code)

		if wsCloseCode, exists := wsCloseCodeMap[appErr.Code]; exists {
			closeCode = wsCloseCode
		}
	}

	logger.WithFields(map[string]interface{}{
		"address":    serverAddress,
		"code":       code,
		"close_code": closeCode,
	}).Error(err)
}
