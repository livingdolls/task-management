package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task-management/internal/applications/services"
	"task-management/internal/config"
	"task-management/internal/infra/adapter/http/handler"
	"task-management/internal/infra/adapter/http/router"
	"task-management/internal/infra/adapter/storages"
	"task-management/internal/infra/logger"
	"task-management/internal/infra/security"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppServer struct {
	DB     *gorm.DB
	Config *config.AppConfig
	Gin    *gin.Engine
}

func InitServer(cf *config.AppConfig, db *gorm.DB) *AppServer {
	engine := gin.Default()

	// Enable CORS
	engine.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	userRepo := storages.NewUserRepository(db)
	jwtService := security.NewJWTAdapter(cf.Secret, time.Hour)
	authService := services.NewAuthService(userRepo, jwtService)
	authHandler := handler.NewAuthHandler(authService)
	taskRepo := storages.NewTaskRepository(db)
	taskService := services.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	// Setup router
	router.SetupRoutes(engine, authHandler, taskHandler, jwtService)

	return &AppServer{
		DB:     db,
		Config: cf,
		Gin:    engine,
	}
}

func StartServer(app *AppServer) *http.Server {
	port := config.Config.Server.Port
	addr := fmt.Sprintf(":%v", port)

	server := &http.Server{
		Addr:    addr,
		Handler: app.Gin,
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
