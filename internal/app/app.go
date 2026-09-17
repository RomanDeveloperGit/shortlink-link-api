package app

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/config"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/logger"
)

func Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()

	logOpts := logger.NewOptions(
		cfg.Env,
		cfg.ServiceName,
		cfg.ServiceVersion,
	)
	log := logger.MustSetup(logOpts)

	httpServer := NewHTTPServer(
		cfg.HTTPServer.Host,
		cfg.HTTPServer.Port,
		cfg.HTTPServer.ReadTimeout,
		cfg.HTTPServer.WriteTimeout,
		cfg.GracefulShutdownPerResourceTimeout,
		log,
	)

	httpServer.BackgroundRun()

	log.Info("server started")

	<-ctx.Done()

	log.Info("shutting down server")

	// Останавливаем слушатель ОС сигналов, чтобы последующий второй CTRL+C мгновенно завершал работу.
	// Это работает, потому что снова произойдет INTERRUPT, на который никто не подписан(т.к. stop убирает эту подписку как раз)
	// А по умолчанию без подписки программа сразу же завершает своё исполнение
	// Сделано на случай, если нужно срочно остановить сервер, не дожидаясь полного graceful shutdown (часто в local-режиме)
	stop()

	httpServer.ShutdownGracefully()
	// postgres.ShutdownGracefully()
	// redis.ShutdownGracefully()
	// kafka.ShutdownGracefully()

	log.Info("server stopped")
}
