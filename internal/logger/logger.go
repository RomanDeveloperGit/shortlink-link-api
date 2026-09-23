package logger

import (
	"log/slog"
	"os"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/config/env"
)

type Options struct {
	File           *os.File
	Env            env.Env
	ServiceName    string
	ServiceVersion string
}

func MustSetup(opts *Options) *slog.Logger {
	var handler slog.Handler

	switch opts.Env {
	case env.EnvLocal:
		handler = newLocalHandler(opts.File)
	case env.EnvDev:
		handler = slog.NewJSONHandler(opts.File, &slog.HandlerOptions{Level: slog.LevelDebug})
	case env.EnvProd:
		handler = slog.NewJSONHandler(opts.File, &slog.HandlerOptions{Level: slog.LevelInfo})
	default:
		// Специально паникуем, чтобы при добавлении нового окружения не забыли настроить логгер осознанно + "защита от дурака" (чтобы не передали фигню в виде строки)
		// Также будет паника, если мы попробуем из env.EnvTest окружения этот логгер инициализировать - пусть свой моковый логгер делает
		panic("env is not supported")
	}

	return slog.New(handler).With(
		slog.String("service", opts.ServiceName),
		slog.String("version", opts.ServiceVersion),
	)
}
