package router

import (
	"task-management/internal/applications/ports/services"
	"task-management/internal/infra/adapter/http/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, authHandler *handler.AuthHandler, jwtSvc services.JWTService) {
	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	// protected := api.Group("/")
	// protected.Use(middleware.JWTAuthMiddleware(jwtSvc))
	// {
	// 	protected.GET("/profile", profileHandler.GetProfile)
	// }
}
