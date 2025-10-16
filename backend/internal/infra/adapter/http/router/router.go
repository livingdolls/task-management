package router

import (
	"task-management/internal/applications/ports/services"
	"task-management/internal/infra/adapter/http/handler"
	"task-management/internal/infra/adapter/http/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, authHandler *handler.AuthHandler, taskHandler *handler.TaskHandler, jwtSvc services.JWTService) {
	api := r.Group("/api/v1")

	// --- Auth Routes ---
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	// --- Protected Routes ---
	protectedGroup := api.Group("/")
	protectedGroup.Use(middleware.JWTMiddleware(jwtSvc))
	{
		// User profile
		protectedGroup.GET("/profile", authHandler.Me)

		// Task routes
		taskGroup := protectedGroup.Group("/tasks")
		{
			taskGroup.POST("/", taskHandler.Create)
			taskGroup.GET("/", taskHandler.Get)
			taskGroup.GET("/:id", taskHandler.GetByID)
			taskGroup.PUT("/:id", taskHandler.Update)
			taskGroup.DELETE("/:id", taskHandler.Delete)
		}
	}
}
