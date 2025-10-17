package handler

import (
	"net/http"
	"task-management/internal/applications/dto/request"
	"task-management/internal/applications/dto/response"
	"task-management/internal/applications/ports/services"
	"task-management/internal/infra/adapter/http/middleware"
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

// Register godoc
// @Summary Register a new user
// @Description Register a new user with name, username, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.RegisterUser true "User registration data"
// @Success 201 {object} response.BaseUserResponse "User registered successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request - invalid input"
// @Failure 409 {object} response.ErrorResponse "Conflict - username already exists"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterUser

	if err := c.ShouldBindJSON(&req); err != nil {
		resp := response.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Error:   err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	user, err := h.auth.Register(req.Name, req.Username, req.Password)

	if err != nil {
		if err.Error() == "username is already exists" {
			resp := response.ErrorResponse{
				Success: false,
				Code:    http.StatusConflict,
				Error:   "Username already exists",
			}
			c.JSON(http.StatusConflict, resp)
			return
		}

		resp := response.ErrorResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Error:   "Internal server error",
		}
		c.JSON(http.StatusInternalServerError, resp)

		logger.Info("failed to register user: ", zap.Error(err))
		return
	}

	resp := response.BaseUserResponse{
		Success: true,
		Code:    http.StatusCreated,
		Data: response.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
		},
	}

	c.JSON(http.StatusCreated, resp)
}

// Login godoc
// @Summary User login
// @Description Authenticate user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.LoginUser true "Login credentials"
// @Success 200 {object} response.BaseAuthResponse "success: true, code: 200, data: response.AuthResponse"
// @Failure 400 {object} response.ErrorResponse "success: false, code: 400, error: validation error"
// @Failure 401 {object} response.ErrorResponse "success: false, code: 401, error: Invalid username or password"
// @Failure 404 {object} response.ErrorResponse "success: false, code: 404, error: User not found"
// @Failure 500 {object} response.ErrorResponse "success: false, code: 500, error: Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginUser

	if err := c.ShouldBindJSON(&req); err != nil {
		resp := response.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Error:   err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	token, user, err := h.auth.Login(req.Username, req.Password)

	if err != nil {
		if err.Error() == "invalid username or password" {
			resp := response.ErrorResponse{
				Success: false,
				Code:    http.StatusUnauthorized,
				Error:   "Invalid username or password",
			}
			c.JSON(http.StatusUnauthorized, resp)
			return
		}

		if err.Error() == "user not found" {
			resp := response.ErrorResponse{
				Success: false,
				Code:    http.StatusNotFound,
				Error:   "User not found",
			}
			c.JSON(http.StatusNotFound, resp)
			return
		}

		resp := response.ErrorResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Error:   "Internal server error",
		}

		c.JSON(http.StatusInternalServerError, resp)
		logger.Info("failed to login user: ", zap.Error(err))
		return
	}

	resp := response.BaseAuthResponse{
		Success: true,
		Code:    http.StatusOK,
		Data: response.AuthResponse{
			Token: token,
			User: &response.UserResponse{
				ID:       user.ID,
				Name:     user.Name,
				Username: user.Username,
			},
		},
	}

	c.JSON(http.StatusOK, resp)
}

// Me godoc
// @Summary Get current user profile
// @Description Get the profile information of the currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "success: true, code: 200, data: UserResponse"
// @Failure 401 {object} response.ErrorResponse "success: false, code: 401, error: Unauthorized"
// @Failure 500 {object} response.ErrorResponse "error: Internal server error"
// @Router /auth/profile [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userClaims, ok := middleware.GetUserClaims(c)

	if !ok {
		resp := response.ErrorResponse{
			Success: false,
			Code:    http.StatusUnauthorized,
			Error:   "Unauthorized",
		}
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	userID := userClaims.UserID

	user, err := h.auth.Me(userID)
	if err != nil {
		resp := response.ErrorResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Error:   "Internal server error",
		}

		logger.Info("failed to get user profile: ", zap.Error(err))
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := response.BaseUserResponse{
		Success: true,
		Code:    http.StatusOK,
		Data: response.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
		},
	}

	c.JSON(http.StatusOK, resp)
}
