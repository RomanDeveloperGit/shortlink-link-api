package config

import (
	"time"

	"github.com/caarlos0/env/v11"

	internalEnv "github.com/RomanDeveloperGit/shortlink-link-api/internal/config/env"
)

type Config struct {
	Env                                internalEnv.Env `env:"ENV,required,notEmpty"`
	ServiceName                        string          `env:"SERVICE_NAME,required,notEmpty"`
	ServiceVersion                     string          `env:"SERVICE_VERSION,required,notEmpty"`
	ShortCodeLength                    int             `env:"SHORT_CODE_LENGTH,required,notEmpty"`
	AttemptsGenerateShortLinkLimit     int             `env:"ATTEMPTS_GENERATE_SHORT_CODE_LIMIT,required,notEmpty"`
	GracefulShutdownPerResourceTimeout time.Duration   `env:"GRACEFUL_SHUTDOWN_PER_RESOURCE_TIMEOUT,required,notEmpty"`
	HTTPServer                         HTTPServer
	PostgreSQL                         PostgreSQL
	Kafka                              Kafka
}

type HTTPServer struct {
	Host         string        `env:"HTTP_HOST,required,notEmpty"`
	Port         int           `env:"HTTP_PORT,required,notEmpty"`
	IdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT,required,notEmpty"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT,required,notEmpty"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT,required,notEmpty"`
}

type PostgreSQL struct {
	Host            string        `env:"POSTGRES_HOST,required,notEmpty"`
	Port            int           `env:"POSTGRES_PORT,required,notEmpty"`
	User            string        `env:"POSTGRES_USER,required,notEmpty"`
	Password        string        `env:"POSTGRES_PASSWORD,required,notEmpty"`
	DB              string        `env:"POSTGRES_DB,required,notEmpty"`
	SSLMode         string        `env:"POSTGRES_SSL_MODE,required,notEmpty"`
	MaxIdleConns    int           `env:"POSTGRES_MAX_IDLE_CONNS,required,notEmpty"`
	MaxOpenConns    int           `env:"POSTGRES_MAX_OPEN_CONNS,required,notEmpty"`
	ConnMaxLifetime time.Duration `env:"POSTGRES_CONN_MAX_LIFETIME,required,notEmpty"`
	PingTimeout     time.Duration `env:"POSTGRES_PING_TIMEOUT,required,notEmpty"`
}

type Kafka struct {
	Brokers  []string `env:"KAFKA_BROKERS,required,notEmpty"`
	User     string   `env:"KAFKA_USER,required,notEmpty"`
	Password string   `env:"KAFKA_PASSWORD,required,notEmpty"`
}

func MustLoad() *Config {
	cfg := &Config{}

	// Специально через Must, чтобы вылезла паника в случае некорректного конфига
	cfg = env.Must(cfg, env.Parse(cfg))

	if !internalEnv.IsEnv(string(cfg.Env)) {
		panic("failed to load config: env is not supported")
	}

	return cfg
}
