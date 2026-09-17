package config

import (
	"time"

	"github.com/caarlos0/env/v11"

	internalEnv "github.com/RomanDeveloperGit/shortlink-link-api/internal/env"
)

type Config struct {
	Env                                internalEnv.Env `env:"ENV,required,notEmpty"`
	ServiceName                        string          `env:"SERVICE_NAME,required,notEmpty"`
	ServiceVersion                     string          `env:"SERVICE_VERSION,required,notEmpty"`
	GracefulShutdownPerResourceTimeout time.Duration   `env:"GRACEFUL_SHUTDOWN_PER_RESOURCE_TIMEOUT,required,notEmpty"`
	HTTPServer                         HTTPServer
	PostgreSQL                         PostgreSQL
}

type HTTPServer struct {
	Host         string        `env:"HTTP_HOST,required,notEmpty"`
	Port         int           `env:"HTTP_PORT,required,notEmpty"`
	IdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT,required,notEmpty"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT,required,notEmpty"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT,required,notEmpty"`
}

type PostgreSQL struct {
	Host     string `env:"POSTGRES_HOST,required,notEmpty"`
	Port     int    `env:"POSTGRES_PORT,required,notEmpty"`
	User     string `env:"POSTGRES_USER,required,notEmpty"`
	Password string `env:"POSTGRES_PASSWORD,required,notEmpty"`
	DB       string `env:"POSTGRES_DB,required,notEmpty"`
}

func MustLoad() *Config {
	cfg := &Config{}

	// Специально через Must, чтобы вылезла паника в случае некорректного конфига
	cfg = env.Must(cfg, env.Parse(cfg))

	if !internalEnv.IsValid(string(cfg.Env)) {
		panic("env is not supported")
	}

	return cfg
}
