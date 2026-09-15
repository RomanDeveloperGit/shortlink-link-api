package config

import (
	"slices"

	"github.com/caarlos0/env/v11"
)

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

var envs = []string{
	string(EnvLocal),
	string(EnvDev),
	string(EnvProd),
}

type HTTPServer struct {
	Address      string `env:"HTTP_ADDRESS,required,notEmpty"`
	IdleTimeout  string `env:"HTTP_IDLE_TIMEOUT,required,notEmpty"`
	ReadTimeout  string `env:"HTTP_READ_TIMEOUT,required,notEmpty"`
	WriteTimeout string `env:"HTTP_WRITE_TIMEOUT,required,notEmpty"`
}

type PostgreSQL struct {
	Host     string `env:"POSTGRES_HOST,required,notEmpty"`
	Port     string `env:"POSTGRES_PORT,required,notEmpty"`
	User     string `env:"POSTGRES_USER,required,notEmpty"`
	Password string `env:"POSTGRES_PASSWORD,required,notEmpty"`
	DB       string `env:"POSTGRES_DB,required,notEmpty"`
}

type Config struct {
	Env            Env    `env:"ENV,required,notEmpty"`
	ServiceName    string `env:"SERVICE_NAME,required,notEmpty"`
	ServiceVersion string `env:"SERVICE_VERSION,required,notEmpty"`
	HTTPServer     HTTPServer
	PostgreSQL     PostgreSQL
}

func MustLoad() *Config {
	var cfg Config

	cfg = env.Must(cfg, env.Parse(&cfg))

	if !isEnvValid(string(cfg.Env)) {
		panic("env is not supported")
	}

	return &cfg
}

func isEnvValid(env string) bool {
	if slices.Contains(envs, env) {
		return true
	}

	return false
}
