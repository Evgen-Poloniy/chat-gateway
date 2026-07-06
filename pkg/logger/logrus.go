package logger

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Evgen-Poloniy/chat-gateway/internal/config"

	"github.com/sirupsen/logrus"
)

var _files []*os.File

func NewLogrusLogger(cfg *config.LoggerConfig) (*logrus.Logger, error) {
	l := logrus.New()

	// Set logger format.
	switch cfg.Format {
	case config.TextFormat:
		l.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default:
		l.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// Set logger level.
	switch cfg.Level {
	case config.TraceLevel:
		l.SetLevel(logrus.TraceLevel)
	case config.DebugLevel:
		l.SetLevel(logrus.DebugLevel)
	case config.InfoLevel:
		l.SetLevel(logrus.InfoLevel)
	case config.WarnLevel:
		l.SetLevel(logrus.WarnLevel)
	case config.ErrorLevel:
		l.SetLevel(logrus.ErrorLevel)
	case config.PanicLevel:
		l.SetLevel(logrus.PanicLevel)
	case config.FatalLevel:
		l.SetLevel(logrus.FatalLevel)
	default:
		l.SetLevel(logrus.InfoLevel)
	}

	// Set logger on output stream (Stdout, Stderr, Files).
	writers := make([]io.Writer, 0, len(cfg.Files))
	_files = make([]*os.File, 0, len(cfg.Files))

	switch cfg.Output {
	case config.Stdout:
		writers = append(writers, os.Stdout)
	case config.Stderr:
		writers = append(writers, os.Stderr)
	default:
		writers = append(writers, os.Stdout)
	}

	for _, filename := range cfg.Files {
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
		}
		writers = append(writers, file)
		_files = append(_files, file)
	}

	l.SetOutput(io.MultiWriter(writers...))

	return l, nil
}

// Closing logger for graceful shutdown.
func Close() error {
	errs := make([]error, 0, len(_files))

	for _, file := range _files {
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	} else if len(errs) == 1 {
		return fmt.Errorf("error of closing files: %w", errs[0])
	}

	return fmt.Errorf("error of closing files: %w", errors.Join(errs...))
}
