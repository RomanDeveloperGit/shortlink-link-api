package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/config"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/infra/broker"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/infra/db"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/logger"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/model"
	linkRepository "github.com/RomanDeveloperGit/shortlink-link-api/internal/repository/link"
	linkService "github.com/RomanDeveloperGit/shortlink-link-api/internal/service/link"
)

const datasetCount = 10000
const workerCount = 20 // подобрано 

func main() {
	cfg := config.MustLoad()
	log := logger.MustSetup(&logger.Options{
		File:           os.Stdout,
		Env:            cfg.Env,
		ServiceName:    cfg.ServiceName,
		ServiceVersion: cfg.ServiceVersion,
	})

	database := db.MustConnect(&db.Options{
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

	brk := broker.MustConnect(&broker.Options{
		Addrs:                   cfg.Kafka.Brokers,
		User:                    cfg.Kafka.User,
		Password:                cfg.Kafka.Password,
		Logger:                  log,
		GracefulShutdownTimeout: cfg.GracefulShutdownPerResourceTimeout,
	})

	linkRepo := linkRepository.NewRepository(database.DB)
	linkSvc := linkService.NewService(
		linkRepo,
		brk.Producer,
		cfg.ShortCodeLength,
		cfg.AttemptsGenerateShortLinkLimit,
	)

	log.Info("server started")

	start := time.Now()

	type input struct {
		id      int
		fullUrl string
		ttlDays int
	}

	type output struct {
		id   int
		link *model.Link
		err  error
	}

	inputCh := make(chan input, datasetCount)
	outputCh := make(chan output, datasetCount)

	commonError := strings.Builder{}
	errCount := 0

	go func() {
		wg := sync.WaitGroup{}

		log.Debug("Создаем Worker'ов",
			slog.Int("count", workerCount),
		)

		wg.Add(workerCount)
		for i := range workerCount {
			log.Debug("Запускаем Worker'а пахать",
				slog.Int("id", i),
			)

			go func() {
				defer func() {
					log.Debug("Worker закончил свою работу",
						slog.Int("id", i),
					)
					wg.Done()
				}()

				for v := range inputCh {
					link, err := linkSvc.Create(context.Background(), v.fullUrl, v.ttlDays)

					outputCh <- output{
						id:   v.id,
						link: link,
						err:  err,
					}
				}
			}()
		}

		log.Debug("Все Worker'ы запущены")

		wg.Wait()

		close(outputCh)

		log.Debug("Worker'ы закончили свою работу, канал output закрыт")
	}()

	go func() {
		log.Debug("Отправляем задачи на исполнение в WorkerPool",
			slog.Int("count", datasetCount),
		)

		for i := range datasetCount {
			inputCh <- input{
				id:      i,
				fullUrl: fmt.Sprintf("https://example.com/%d", i),
				ttlDays: 1,
			}
		}

		close(inputCh)

		log.Debug("Все задачи на исполнение отправлены, канал input закрыт")
	}()

	log.Debug("Начинаем считывать ответы и формировать итоговый общий")

	for v := range outputCh {
		if v.err != nil {
			errCount++

			fmt.Fprintf(&commonError, "%d: %s\n", v.id, v.err.Error())
		}
	}

	log.Info("loadtest dataset filled",
		slog.Duration("duration", time.Since(start)),
		slog.Any("error", fmt.Errorf("failed to fill loadtest dataset has %d errors: %s", errCount, commonError.String())),
	)

	log.Info("shutting down server")

	database.ShutdownGracefully()
	brk.ShutdownGracefully()
	// redis.ShutdownGracefully()

	log.Info("server stopped")
}
