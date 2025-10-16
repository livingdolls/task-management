package handler

import (
	"net/http"
	"task-management/internal/applications/dto/request"
	"task-management/internal/applications/dto/response"
	"task-management/internal/applications/ports/services"
	"task-management/internal/infra/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	auth services.AuthService
}

func NewAuthHandler(auth services.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterUser

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"error":   err.Error(),
		})
		return
	}

	user, err := h.auth.Register(req.Name, req.Username, req.Password)

	if err != nil {
		if err.Error() == "username is already exists" {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"code":    http.StatusConflict,
				"error":   "Username already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})

		logger.Info("failed to register user: ", zap.Error(err))
		return
	}

	resp := response.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"code":    http.StatusCreated,
		"data":    resp,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginUser

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"error":   err.Error(),
		})
		return
	}

	token, user, err := h.auth.Login(req.Username, req.Password)

	if err != nil {
		if err.Error() == "invalid username or password" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"error":   "Invalid username or password",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})

		logger.Info("failed to login user: ", zap.Error(err))
		return
	}

	resp := response.AuthResponse{
		Token: token,
		User: &response.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"data":    resp,
	})
}
