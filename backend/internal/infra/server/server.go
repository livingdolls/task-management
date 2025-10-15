package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task-management/internal/config"
	"task-management/internal/infra/logger"
	"time"

	"go.uber.org/zap"
)

func StartServer() *http.Server {
	port := config.Config.Server.Port
	addr := fmt.Sprintf(":%v", port)

	server := &http.Server{
		Addr: addr,
	}

	logger.Info("Starting server", zap.String("address", addr))

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server start failed", zap.Error(err))
		}
	}()

	return server
}

// WaitForShutdown menunggu sinyal OS dan melakukan graceful shutdown.
func WaitForShutdown(server *http.Server, cleanupFuncs ...func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	)

	sig := <-quit
	logger.Info("Received shutdown signal", zap.String("signal", sig.String()))

	for _, fn := range cleanupFuncs {
		fn()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	} else {
		logger.Info("Server shutdown gracefully")
	}

	signal.Stop(quit)
	close(quit)
}
