package cache

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/infra/tryshutdowngracefully"
)

type cache struct {
	RDB                     *redis.Client
	logger                  *slog.Logger
	gracefulShutdownTimeout time.Duration
}

type Options struct {
	Host                    string
	Port                    int
	User                    string
	Password                string
	PingTimeout             time.Duration
	Logger                  *slog.Logger
	GracefulShutdownTimeout time.Duration
}

func MustConnect(opts *Options) *cache {
	addr := fmt.Sprintf("%s:%d", opts.Host, opts.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: opts.Password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), opts.PingTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		panic(fmt.Errorf("failed to ping redis: %v", err))
	}

	opts.Logger.Debug("redis connected",
		slog.String("host", opts.Host),
		slog.Int("port", opts.Port),
		slog.String("user", opts.User),
	)

	return &cache{
		RDB:                     client,
		logger:                  opts.Logger,
		gracefulShutdownTimeout: opts.GracefulShutdownTimeout,
	}
}

func (c *cache) TryShutdownGracefully() {
	tryshutdowngracefully.TryShutdownGracefully(
		func(done chan<- error) {
			done <- c.RDB.Close()
		},
		"redis",
		c.gracefulShutdownTimeout,
		c.logger,
	)
}
