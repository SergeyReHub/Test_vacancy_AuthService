package gracefull_shutdown

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
)

func GracefulShutDown(logger *zap.Logger, servers ...*http.Server) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Ожидаем сигнал завершения
	<-signalChan
	log.Println("Received shutdown signal, shutting down gracefully...")

	var wg sync.WaitGroup
	for _, srv := range servers {
		if srv != nil {
			wg.Add(1)
			go func(s *http.Server) {
				defer wg.Done()
				err := s.Shutdown(context.Background())
				if err != nil {
					logger.Error("Error occured than shutdown server")
				}
				log.Println("server stopped")
			}(srv)
		}
	}

	wg.Wait()
	log.Println("All servers stopped")
}
