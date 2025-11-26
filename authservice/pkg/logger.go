package pkg

import (
	"log/slog"
	"os"
)

func New(level slog.Level) *slog.Logger {
	var handler slog.Handler

	handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	logger := slog.New(handler)

	logger.Info("logger initialized",
		slog.String("level", level.String()),
	)

	return logger
}
