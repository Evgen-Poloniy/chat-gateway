package logger

import (
	"github.com/Evgen-Poloniy/chat-gateway/internal/config"

	"github.com/sirupsen/logrus"
)

func NewLogrusLogger(cfg *config.LoggerConfig) *logrus.Logger {
	logger := logrus.New()

	// Set logger format.
	switch cfg.Format {
	case config.TextFormat:
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default:
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// Set logger level.
	switch cfg.Level {
	case config.TraceLevel:
		logger.SetLevel(logrus.TraceLevel)
	case config.DebugLevel:
		logger.SetLevel(logrus.DebugLevel)
	case config.InfoLevel:
		logger.SetLevel(logrus.InfoLevel)
	case config.WarnLevel:
		logger.SetLevel(logrus.WarnLevel)
	case config.ErrorLevel:
		logger.SetLevel(logrus.ErrorLevel)
	case config.PanicLevel:
		logger.SetLevel(logrus.PanicLevel)
	case config.FatalLevel:
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger
}
