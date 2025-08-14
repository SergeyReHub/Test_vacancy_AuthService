package main

import (
	"auth/internal/config"
	"auth/internal/server"
	"auth/pkg/gracefull_shutdown"
	"auth/pkg/logger"
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
