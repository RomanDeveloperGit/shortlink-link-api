package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/config"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/logger"
)

func main() {
	// Заданный порядок инициализации важен, чтобы не словить панику
	config.MustLoad()
	logger.MustSetup()

	cfg := config.Get()
	log := logger.Get()

	mux := http.NewServeMux()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port),
		Handler:      mux,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", slog.String("error", err.Error()))
		}
	}()

	log.Info("server started")

	<-ctx.Done()

	log.Info("shutting down server")

	// Останавливаем слушатель ОС сигналов, чтобы последующий второй CTRL+C мгновенно завершал работу.
	// Это работает, потому что снова произойдет INTERRUPT, на который никто не подписан(т.к. stop убирает эту подписку как раз)
	// А по умолчанию без подписки программа сразу же завершает своё исполнение
	// Сделано на случай, если нужно срочно остановить сервер, не дожидаясь полного graceful shutdown (часто в local-режиме)
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	log.Debug("HTTP server closed")

	time.Sleep(3 * time.Second)

	log.Debug("DB, Redis, etc. closed")

	log.Info("server stopped")
}
