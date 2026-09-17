package env

import "slices"

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

func IsValid(env string) bool {
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
