package server

import (
	"auth/internal/config"
	"auth/internal/service"
	http_service "auth/internal/transport/http"
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func StartServer(cfg *config.Config, logger *zap.Logger) (*http.Server, error) {
	src, err := service.NewAuthService(logger, cfg)
	if err != nil {
		logger.Error("Error create service.", zap.Error(err))
		return nil, err
	}

	r := http_service.NewRouter(logger, src, context.Background())
	srv := &http.Server{
		Handler:      r,
		Addr:         cfg.HTTPServer.GetAddr(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return srv, nil
}
