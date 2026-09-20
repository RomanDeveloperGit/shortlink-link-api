package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/config"
	healthHandler "github.com/RomanDeveloperGit/shortlink-link-api/internal/handler/health"
	linkHandler "github.com/RomanDeveloperGit/shortlink-link-api/internal/handler/link"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/infra/db"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/logger"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/middleware"
	linkRepository "github.com/RomanDeveloperGit/shortlink-link-api/internal/repository/link"
	linkService "github.com/RomanDeveloperGit/shortlink-link-api/internal/service/link"
	httpserver "github.com/RomanDeveloperGit/shortlink-link-api/internal/transport/httpserver"
	"github.com/go-playground/validator/v10"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()
	log := logger.MustSetup(&logger.Options{
		File:           os.Stdout,
		Env:            cfg.Env,
		ServiceName:    cfg.ServiceName,
		ServiceVersion: cfg.ServiceVersion,
	})

	database := db.MustConnect(&db.Config{
		Host:                    cfg.PostgreSQL.Host,
		Port:                    cfg.PostgreSQL.Port,
		User:                    cfg.PostgreSQL.User,
		Password:                cfg.PostgreSQL.Password,
		DBName:                  cfg.PostgreSQL.DB,
		SSLMode:                 cfg.PostgreSQL.SSLMode,
		MaxIdleConns:            cfg.PostgreSQL.MaxIdleConns,
		MaxOpenConns:            cfg.PostgreSQL.MaxOpenConns,
		ConnMaxLifetime:         cfg.PostgreSQL.ConnMaxLifetime,
		PingTimeout:             cfg.PostgreSQL.PingTimeout,
		GracefulShutdownTimeout: cfg.GracefulShutdownPerResourceTimeout,
		Logger:                  log,
	})

	validator := validator.New()

	healthH := healthHandler.NewHandler()

	linkRepo := linkRepository.NewRepository(database.DB)
	linkSvc := linkService.NewService(linkRepo, log, cfg.ShortLinkLength, cfg.AttemptsGenerateShortLinkLimit)
	linkH := linkHandler.NewHandler(linkSvc, validator)

	mid := middleware.Wrapper(log)

	httpServer := httpserver.NewHTTPServer(&httpserver.Options{
		Host:                 cfg.HTTPServer.Host,
		Port:                 cfg.HTTPServer.Port,
		ReadTimeout:          cfg.HTTPServer.ReadTimeout,
		WriteTimeout:         cfg.HTTPServer.WriteTimeout,
		GracefulShutdownTime: cfg.GracefulShutdownPerResourceTimeout,
		Logger:               log,
		Middleware:           mid,
		Handlers: httpserver.Handlers{
			HealthHandler: healthH,
			LinkHandler:   linkH,
		},
	})

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
	database.ShutdownGracefully()
	// redis.ShutdownGracefully()
	// kafka.ShutdownGracefully()

	log.Info("server stopped")
}
