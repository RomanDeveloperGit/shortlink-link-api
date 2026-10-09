package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/config"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/infra/broker"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/infra/cache"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/infra/db"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/logger"
	linkRepository "github.com/RomanDeveloperGit/shortlink-link-api/internal/repository/link"
	linkService "github.com/RomanDeveloperGit/shortlink-link-api/internal/service/link"
	httpserver "github.com/RomanDeveloperGit/shortlink-link-api/internal/transport/httpserver"
	healthHandler "github.com/RomanDeveloperGit/shortlink-link-api/internal/transport/httpserver/handler/health"
	linkHandler "github.com/RomanDeveloperGit/shortlink-link-api/internal/transport/httpserver/handler/link"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/transport/httpserver/middleware"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/transport/httpserver/requestvalidation"
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

	dbClient := db.MustConnect(&db.Options{
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

	brokerClient := broker.MustConnect(&broker.Options{
		ServiceName:             cfg.ServiceName,
		Addrs:                   cfg.Kafka.Brokers,
		User:                    cfg.Kafka.User,
		Password:                cfg.Kafka.Password,
		Logger:                  log,
		GracefulShutdownTimeout: cfg.GracefulShutdownPerResourceTimeout,
	})

	cacheClient := cache.MustConnect(&cache.Options{
		Host:                    cfg.Redis.Host,
		Port:                    cfg.Redis.Port,
		User:                    cfg.Redis.User,
		Password:                cfg.Redis.Password,
		PingTimeout:             cfg.Redis.PingTimeout,
		Logger:                  log,
		GracefulShutdownTimeout: cfg.GracefulShutdownPerResourceTimeout,
	})

	validator := requestvalidation.NewValidator()

	healthH := healthHandler.NewHandler()

	linkRepo := linkRepository.NewRepository(dbClient.DB)
	linkSvc := linkService.NewService(
		linkRepo,
		brokerClient.Producer,
		cacheClient.RDB,
		log,
		cfg.ShortCodeLength,
		cfg.AttemptsGenerateShortLinkLimit,
	)
	linkH := linkHandler.NewHandler(linkSvc, validator, log)

	combinedMiddleware := middleware.Combine(
		middleware.Recovery(log),
		middleware.TraceID(),
		middleware.RequestID(),
		middleware.Logging(log),
	)

	httpServer := httpserver.NewHTTPServer(&httpserver.Options{
		Host:                 cfg.HTTPServer.Host,
		Port:                 cfg.HTTPServer.Port,
		ReadTimeout:          cfg.HTTPServer.ReadTimeout,
		WriteTimeout:         cfg.HTTPServer.WriteTimeout,
		GracefulShutdownTime: cfg.GracefulShutdownPerResourceTimeout,
		Logger:               log,
		Middleware:           combinedMiddleware,
		Handlers: httpserver.Handlers{
			HealthHandler: healthH,
			LinkHandler:   linkH,
		},
	})

	go httpServer.Run()

	time.Sleep(1 * time.Millisecond)

	log.Info("server started")

	<-ctx.Done()

	log.Info("shutting down server")

	// Останавливаем слушатель ОС сигналов, чтобы последующий второй CTRL+C мгновенно завершал работу.
	// Это работает, потому что снова произойдет INTERRUPT, на который никто не подписан(т.к. stop убирает эту подписку как раз)
	// А по умолчанию без подписки программа сразу же завершает своё исполнение
	// Сделано на случай, если нужно срочно остановить сервер, не дожидаясь полного graceful shutdown (часто в local-режиме)
	stop()

	httpServer.ShutdownGracefully()
	cacheClient.TryShutdownGracefully()
	brokerClient.TryShutdownGracefully()
	dbClient.TryShutdownGracefully()

	log.Info("server stopped")
}
