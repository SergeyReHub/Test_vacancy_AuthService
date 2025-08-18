package server

import (
	"auth/backend/internal/config"
	"auth/backend/internal/service"
	http_service "auth/backend/internal/transport/http"
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

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("Server started", zap.String("addr", cfg.HTTPServer.GetAddr()))
		err := srv.ListenAndServe()
		if err != nil {
			serverErr <- err
		}
		close(serverErr)
	}()
	select {
    case err := <-serverErr:
		logger.Error("Error by starting server", zap.Error(err))
        return nil, err
    case <-time.After(100 * time.Millisecond): // brief pause to detect immediate failures
        return srv, nil
    }
}
