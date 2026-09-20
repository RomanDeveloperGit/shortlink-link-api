package env

import "slices"

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
	EnvTest  Env = "test"
)

func IsEnv(str string) bool {
	envs := []string{
		string(EnvLocal),
		string(EnvDev),
		string(EnvProd),
		string(EnvTest),
	}

	if slices.Contains(envs, str) {
		return true
	}

	return false
}
