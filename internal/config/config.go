package config

import (
	"slices"
	"time"

	"github.com/caarlos0/env/v11"
)

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

type Config struct {
	Env                     Env           `env:"ENV,required,notEmpty"`
	ServiceName             string        `env:"SERVICE_NAME,required,notEmpty"`
	ServiceVersion          string        `env:"SERVICE_VERSION,required,notEmpty"`
	GracefulShutdownTimeout time.Duration `env:"GRACEFUL_SHUTDOWN_TIMEOUT,required,notEmpty"`
	HTTPServer              HTTPServer
	PostgreSQL              PostgreSQL
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

var config *Config

func isEnvValid(env string) bool {
	envs := []string{
		string(EnvLocal),
		string(EnvDev),
		string(EnvProd),
	}

	if slices.Contains(envs, env) {
		return true
	}

	return false
}

func MustLoad() {
	if config != nil {
		return
	}

	config = &Config{}
	config = env.Must(config, env.Parse(config))

	if !isEnvValid(string(config.Env)) {
		panic("env is not supported")
	}
}

func Get() *Config {
	if config == nil {
		panic("config: MustLoad() was not called")
	}

	return config
}
