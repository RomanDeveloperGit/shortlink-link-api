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

type Options struct {
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

func MustConnect(opts *Options) *database {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		opts.Host, opts.Port, opts.User, opts.Password, opts.DBName, opts.SSLMode)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %v", err))
	}

	db.SetMaxIdleConns(opts.MaxIdleConns)
	db.SetMaxOpenConns(opts.MaxOpenConns)
	db.SetConnMaxLifetime(opts.ConnMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), opts.PingTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		panic(fmt.Errorf("failed to ping database: %v", err))
	}

	opts.Logger.Debug("database connected",
		slog.String("host", opts.Host),
		slog.Int("port", opts.Port),
		slog.String("user", opts.User),
		slog.String("dbname", opts.DBName),
		slog.String("sslmode", opts.SSLMode),
	)

	return &database{
		DB:                      db,
		logger:                  opts.Logger,
		gracefulShutdownTimeout: opts.GracefulShutdownTimeout,
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
			db.logger.Error("failed to shutdown database", slog.Any("error", err))
			return
		}

		db.logger.Debug("database closed")
	case <-ctx.Done():
		// Но всё, что запустилось при вызове db.DB.Close() может продолжать в фоне висеть и работать
		// Допускаем такое, потому что у нас стадия Graceful Shutdown
		db.logger.Error("database shutdown timed out", slog.Any("error", ctx.Err()))
	}
}
