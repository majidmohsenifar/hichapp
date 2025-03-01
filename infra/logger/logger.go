package logger

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	opt := slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewTextHandler(os.Stdout, &opt)
	return slog.New(handler)
}
