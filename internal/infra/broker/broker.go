package broker

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"

	"github.com/RomanDeveloperGit/shortlink-link-api/internal/infra/tryshutdowngracefully"
)

type broker struct {
	Producer                sarama.SyncProducer
	logger                  *slog.Logger
	gracefulShutdownTimeout time.Duration
}

type Options struct {
	ServiceName             string
	Addrs                   []string
	User                    string
	Password                string
	Logger                  *slog.Logger
	GracefulShutdownTimeout time.Duration
}

func MustConnect(opts *Options) *broker {
	cfg := sarama.NewConfig()

	cfg.ClientID = opts.ServiceName

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

func (b *broker) TryShutdownGracefully() {
	tryshutdowngracefully.TryShutdownGracefully(
		func(done chan<- error) {
			done <- b.Producer.Close()
		},
		"kafka",
		b.gracefulShutdownTimeout,
		b.logger,
	)
}
