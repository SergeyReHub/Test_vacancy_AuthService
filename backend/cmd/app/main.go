package main

import (
	"auth/backend/internal/config"
	"auth/backend/internal/server"
	"auth/backend/pkg/gracefull_shutdown"
	"auth/backend/pkg/logger"
)

func main() {
	cfg := config.GetConfig()

	logger := logger.New(&cfg.Logger)

	server, err := server.StartServer(cfg, logger)
	if err != nil {
		return
	}

	gracefull_shutdown.GracefulShutDown(logger, server)
}
