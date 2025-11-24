// Package logger содержит инициализацию логгера приложения.
package logger

import (
	"log/slog"
	"os"
)

// New возвращает JSON-логгер уровня Info.
func New() *slog.Logger {
	return NewWithLevel(slog.LevelInfo)
}

// NewWithLevel создает логгер с заданным уровнем логирования.
func NewWithLevel(level slog.Level) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return slog.New(handler)
}
