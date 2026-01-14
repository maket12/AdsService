package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	var level slog.Level
	var handler slog.Handler

	level = slog.LevelDebug
	handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	logger := slog.New(handler)

	logger.Info("logger initialized",
		slog.String("level", level.String()),
	)

	return logger
}
