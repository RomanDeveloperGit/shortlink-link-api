package logger

import (
	"log/slog"
	"os"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/env"
)

type options struct {
	env            env.Env
	serviceName    string
	serviceVersion string
}

// Специально сделано, чтобы извне не создавать структуру, где можно забыть указать что-то, а просто дернуть конструктор с контрактом
func NewOptions(e env.Env, serviceName, serviceVersion string) *options {
	return &options{
		env:            e,
		serviceName:    serviceName,
		serviceVersion: serviceVersion,
	}
}

func MustSetup(opts *options) *slog.Logger {
	if opts == nil {
		panic("options is nil")
	}

	var log *slog.Logger

	switch opts.env {
	case env.EnvLocal:
		log = slog.New(
			newLocalHandler(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			),
		)
	case env.EnvDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case env.EnvProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		// Специально паникуем, чтобы при добавлении нового окружения не забыли настроить и логгер осознанно
		// + "защита от дурака"
		panic("env is not supported")
	}

	log = log.With(
		slog.String("service", opts.serviceName),
		slog.String("version", opts.serviceVersion),
	)

	return log
}
