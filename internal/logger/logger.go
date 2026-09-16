package logger

import (
	"log/slog"
	"os"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/config"
)

var logger *slog.Logger

func MustSetup() {
	if logger != nil {
		return
	}

	cfg := config.Get()

	switch cfg.Env {
	case config.EnvLocal:
		logger = slog.New(
			NewLocalHandler(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			),
		)
	case config.EnvDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case config.EnvProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	logger = logger.With(
		slog.String("service", cfg.ServiceName),
		slog.String("version", cfg.ServiceVersion),
	)
}

func Get() *slog.Logger {
	if logger == nil {
		panic("logger: MustSetup() was not called")
	}

	return logger
}
