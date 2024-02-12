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
	case "WARN":
		slogLevel = slog.LevelWarn
	case "DEBUG":
		slogLevel = slog.LevelDebug
	case "ERROR":
		slogLevel = slog.LevelError
	}

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel}))

	return Logger{logger: log}
}

func (l Logger) Info(msg string, attrs ...any) {
	l.logger.Info(msg, attrs...)
}

func (l Logger) Error(msg string, attrs ...any) {
	l.logger.Error(msg, attrs...)
}

func (l Logger) Debug(msg string, attrs ...any) {
	l.logger.Debug(msg, attrs...)
}

func (l Logger) Warn(msg string, attrs ...any) {
	l.logger.Warn(msg, attrs...)
}
