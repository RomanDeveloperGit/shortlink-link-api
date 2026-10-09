package tryshutdowngracefully

import (
	"context"
	"log/slog"
	"time"
)

func TryShutdownGracefully(
	closer func(chan<- error),
	resourceName string,
	gracefulShutdownTimeout time.Duration,
	logger *slog.Logger,
) {
	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	done := make(chan error, 1)

	go func() {
		closer(done)
		close(done)
	}()

	select {
	case err := <-done:
		if err != nil {
			logger.Error(
				"failed to shutdown",
				slog.String("resource", resourceName),
				slog.Any("error", err),
			)
			return
		}

		logger.Debug("resource closed", slog.String("resource", resourceName))
	case <-ctx.Done():
		// Тут возможна проблема: если таймаут выйдет, а горутина closer продолжит висеть, то она так и будет в фоне
		// продолжать работать, но мы уже контролировать ее не будем.
		// Из этой функции мы выйдем, результат done никто считывать не будет
		// В общем: нечистое закрытие ресурса + утечка. Но допустимо в процессе Graceful Shutdown

		logger.Error("kafka producer shutdown timed out", slog.Any("error", ctx.Err()))
	}
}
