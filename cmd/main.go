package main

import (
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/config"
	"github.com/RomanDeveloperGit/shortlink-link-api/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.MustSetup(cfg)

	log.Info("starting server")
}
