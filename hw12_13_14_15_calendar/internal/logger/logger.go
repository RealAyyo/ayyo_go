package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func New(level string) Logger {
	var slogLevel slog.Level
	switch level {
	case "INFO":
		slogLevel = slog.LevelInfo
		break
	case "WARN":
		slogLevel = slog.LevelWarn
		break
	case "DEBUG":
		slogLevel = slog.LevelDebug
		break
	case "ERROR":
		slogLevel = slog.LevelError
		break
	}
	var log *slog.Logger

	log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel}))

	return Logger{logger: log}
}

func (l Logger) Info(msg string, attrs ...any) {
	if len(attrs) > 0 {
		l.logger.Info(msg, attrs)
	} else {
		l.logger.Info(msg)
	}
}

func (l Logger) Error(msg string, attrs ...any) {
	if len(attrs) > 0 {
		l.logger.Error(msg, attrs)
	} else {
		l.logger.Error(msg)
	}
}

func (l Logger) Debug(msg string, attrs ...any) {
	if len(attrs) > 0 {
		l.logger.Debug(msg, attrs...)
	} else {
		l.logger.Debug(msg)
	}
}

func (l Logger) Warn(msg string, attrs ...any) {
	if len(attrs) > 0 {
		l.logger.Warn(msg, attrs)
	} else {
		l.logger.Warn(msg)
	}
}
