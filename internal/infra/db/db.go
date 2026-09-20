package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type database struct {
	DB                      *sqlx.DB
	logger                  *slog.Logger
	gracefulShutdownTimeout time.Duration
}

type Config struct {
	Host                    string
	Port                    int
	User                    string
	Password                string
	DBName                  string
	SSLMode                 string
	MaxIdleConns            int
	MaxOpenConns            int
	ConnMaxLifetime         time.Duration
	PingTimeout             time.Duration
	GracefulShutdownTimeout time.Duration
	Logger                  *slog.Logger
}

func MustConnect(cfg *Config) *database {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Open("postgres", dsn)

	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.PingTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}

	cfg.Logger.Debug("database connected")

	return &database{
		DB:                      db,
		logger:                  cfg.Logger,
		gracefulShutdownTimeout: cfg.GracefulShutdownTimeout,
	}
}

func (db *database) ShutdownGracefully() {
	ctx, cancel := context.WithTimeout(context.Background(), db.gracefulShutdownTimeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- db.DB.Close()
	}()

	select {
	case err := <-done:
		if err != nil {
			db.logger.Error("failed to shutdown database", slog.String("error", err.Error()))
			return
		}

		db.logger.Debug("database closed")
	case <-ctx.Done():
		// Но всё, что запустилось при вызове db.DB.Close() может продолжать в фоне висеть и работать
		// Допускаем такое, потому что у нас стадия Graceful Shutdown
		db.logger.Error("database shutdown timed out", slog.String("error", ctx.Err().Error()))
	}
}
