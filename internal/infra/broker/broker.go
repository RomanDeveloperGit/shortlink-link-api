package broker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
)

type broker struct {
	Producer                sarama.SyncProducer
	logger                  *slog.Logger
	gracefulShutdownTimeout time.Duration
}

type Options struct {
	Addrs                   []string
	User                    string
	Password                string
	Logger                  *slog.Logger
	GracefulShutdownTimeout time.Duration
}

func MustConnect(opts *Options) *broker {
	cfg := sarama.NewConfig()

	cfg.ClientID = opts.User

	// cfg.Net.SASL.Enable = true
	// cfg.Net.SASL.User = opts.User
	// cfg.Net.SASL.Password = opts.Password
	// cfg.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	// cfg.Net.SASL.SCRAMClientGeneratorFunc = SCRAMClientGeneratorFunc

	cfg.Producer.Retry.Max = 3
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(opts.Addrs, cfg)
	if err != nil {
		panic(fmt.Errorf("failed to create kafka producer: %v", err))
	}

	opts.Logger.Debug("kafka producer connected",
		slog.String("addrs", fmt.Sprintf("%v", opts.Addrs)),
		slog.String("user", opts.User),
	)

	return &broker{
		Producer:                producer,
		logger:                  opts.Logger,
		gracefulShutdownTimeout: opts.GracefulShutdownTimeout,
	}
}

func (b *broker) ShutdownGracefully() {
	ctx, cancel := context.WithTimeout(context.Background(), b.gracefulShutdownTimeout)
	defer cancel()

	done := make(chan error, 1)

	go func() {
		done <- b.Producer.Close()
	}()

	select {
	case err := <-done:
		if err != nil {
			b.logger.Error("failed to shutdown kafka producer", slog.Any("error", err))
			return
		}

		b.logger.Debug("kafka producer closed")
	case <-ctx.Done():
		// аналогичная проблема утечки и нечистого закрытия, как в БД.
		// Допускаем это

		b.logger.Error("kafka producer shutdown timed out", slog.Any("error", ctx.Err()))
	}
}
