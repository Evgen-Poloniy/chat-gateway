package ws

import (
	"github.com/sirupsen/logrus"
)

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
