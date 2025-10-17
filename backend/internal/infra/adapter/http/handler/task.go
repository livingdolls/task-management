package handler

import (
	"net/http"
	"strconv"
	"task-management/internal/applications/dto/request"
	"task-management/internal/applications/dto/response"
	"task-management/internal/applications/ports/services"
	"task-management/internal/domain"
	"task-management/internal/infra/adapter/http/middleware"
	"task-management/internal/infra/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TaskHandler struct {
	taskService services.TaskService
}

func NewTaskHandler(taskService services.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

// Create creates a new task for the authenticated user
// @Summary Create a new task
// @Description Create a new task with the provided details for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body request.CreateTask true "Task creation request"
// @Success 201 {object} response.BaseTaskResponse "Task created successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request - invalid JSON or validation error"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - invalid or missing token"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
	var req request.CreateTask

	if err := c.ShouldBindJSON(&req); err != nil {
		resp := response.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Error:   err.Error(),
		}

		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// claims token dari middleware
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

	task := domain.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Deadline:    req.Deadline,
	}

	if err := h.taskService.CreateTask(userClaims.UserID, &task); err != nil {
		resp := response.ErrorResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Error:   "Internal server error",
		}

		c.JSON(http.StatusInternalServerError, resp)
		logger.Info("failed to create task: ", zap.Error(err))
		return
	}

	resp := response.BaseTaskResponse{
		Success: true,
		Code:    http.StatusCreated,
		Data: response.Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      string(task.Status),
			Deadline:    task.Deadline,
			CreatedAt:   task.CreatedAt,
		},
	}

	c.JSON(http.StatusCreated, resp)
}

// Get retrieves tasks for the authenticated user with optional filtering
// @Summary Get user tasks
// @Description Retrieves a list of tasks for the authenticated user with optional filtering by status and deadline
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by task status" Enums(pending, in_progress, completed)
// @Param deadline query string false "Filter by deadline date (YYYY-MM-DD format)" Format(date)
// @Success 200 {object} response.ListTaskResponse "Successfully retrieved tasks"
// @Failure 400 {object} response.ErrorResponse "Invalid deadline format"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - invalid or missing token"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /tasks [get]
func (h *TaskHandler) Get(c *gin.Context) {
	// claims token dari middleware
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

	var status *domain.TaskStatus
	if s := c.Query("status"); s != "" {
		ts := domain.TaskStatus(s)
		status = &ts
	}

	var deadline *time.Time

	if d := c.Query("deadline"); d != "" {
		parsedDeadline, err := time.Parse("2006-01-02", d)
		if err != nil {
			resp := response.ErrorResponse{
				Success: false,
				Code:    http.StatusBadRequest,
				Error:   "Invalid deadline format. Use YYYY-MM-DD.",
			}

			c.JSON(http.StatusBadRequest, resp)
			return
		}
		deadline = &parsedDeadline
	}

	tasks, err := h.taskService.GetTasks(userClaims.UserID, status, deadline)

	if err != nil {
		resp := response.ErrorResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Error:   "Internal server error",
		}

		c.JSON(http.StatusInternalServerError, resp)

		logger.Info("failed to get tasks: ", zap.Error(err))
		return
	}

	var resp []response.Task

	for _, task := range tasks {
		resp = append(resp, response.Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      string(task.Status),
			Deadline:    task.Deadline,
			CreatedAt:   task.CreatedAt,
		})
	}

	response := response.ListTaskResponse{
		Success: true,
		Code:    http.StatusOK,
		Data:    resp,
	}

	c.JSON(http.StatusOK, response)

}

// GetByID godoc
// @Summary Get task by ID
// @Description Retrieve a specific task by its ID for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Security BearerAuth
// @Success 200 {object} response.BaseTaskResponse "Task retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid task ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Task not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		resp := response.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Error:   "Invalid task ID",
		}

		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// claims token dari middleware
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

	task, err := h.taskService.GetTaskById(uint(id), userClaims.UserID)

	if err != nil {
		if err.Error() == "unauthorized" {
			resp := response.ErrorResponse{
				Success: false,
				Code:    http.StatusUnauthorized,
				Error:   "Unauthorized",
			}

			c.JSON(http.StatusUnauthorized, resp)
			return
		}

		if err.Error() == "task not found" {
			resp := response.ErrorResponse{
				Success: false,
				Code:    http.StatusNotFound,
				Error:   "Task not found",
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

		logger.Info("failed to get task by id: ", zap.Error(err))
		return
	}

	if task == nil {
		resp := response.ErrorResponse{
			Success: false,
			Code:    http.StatusNotFound,
			Error:   "Task not found",
		}

		c.JSON(http.StatusNotFound, resp)
		return
	}

	response := response.BaseTaskResponse{
		Success: true,
		Code:    http.StatusOK,
		Data: response.Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      string(task.Status),
			Deadline:    task.Deadline,
			CreatedAt:   task.CreatedAt,
		},
	}

	c.JSON(http.StatusOK, response)
}

// Update updates an existing task for the authenticated user
// @Summary Update an existing task
// @Description Update a task with the provided details for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body request.UpdateTask true "Task update request"
// @Success 200 {object} response.BaseTaskResponse "Task updated successfully"
// @Failure 400 {object} response.ErrorResponse "Bad request - invalid JSON, validation error, or invalid task ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - invalid or missing token"
// @Failure 404 {object} response.ErrorResponse "Task not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /tasks/{id} [put]
func (h *TaskHandler) Update(c *gin.Context) {
	var req request.UpdateTask

	// Get task ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		resp := response.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Error:   "Invalid task ID",
		}

		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// claims token dari middleware
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

	if err := c.ShouldBindJSON(&req); err != nil {
		resp := response.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Error:   err.Error(),
		}

		c.JSON(http.StatusBadRequest, resp)
		return
	}

	task := domain.Task{
		ID:          uint(id),
		Title:       derefString(req.Title),
		Description: derefString(req.Description),
		Status:      derefStatus(req.Status),
		Deadline:    req.Deadline,
	}

	if err := h.taskService.UpdateTask(&task, userClaims.UserID); err != nil {
		if err.Error() == "unauthorized" {
			resp := response.ErrorResponse{
				Success: false,
				Code:    http.StatusUnauthorized,
				Error:   "Unauthorized",
			}

			c.JSON(http.StatusUnauthorized, resp)
			return
		}

		if err.Error() == "task not found" {
			resp := response.ErrorResponse{
				Success: false,
				Code:    http.StatusNotFound,
				Error:   "Task not found",
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

		logger.Error("failed to update task: ", zap.Error(err))
		return
	}

	resp := response.BaseTaskResponse{
		Success: true,
		Code:    http.StatusOK,
		Data: response.Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      string(task.Status),
			Deadline:    task.Deadline,
			CreatedAt:   task.CreatedAt,
		},
	}

	c.JSON(http.StatusOK, resp)
}

// Delete godoc
// @Summary Delete a task
// @Description Delete a task by ID. Only the task owner can delete their own task.
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Security BearerAuth
// @Success 200 {object} response.DeleteResponse "Task deleted successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid task ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Task not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		resp := response.ErrorResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Error:   "Invalid task ID",
		}

		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// claims token dari middleware
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

	if err := h.taskService.DeleteTask(uint(id), userClaims.UserID); err != nil {
		if err.Error() == "unauthorized" {
			resp := response.ErrorResponse{
				Success: false,
				Code:    http.StatusUnauthorized,
				Error:   "Unauthorized",
			}

			c.JSON(http.StatusUnauthorized, resp)
			return
		}

		if err.Error() == "task not found" {
			resp := response.ErrorResponse{
				Success: false,
				Code:    http.StatusNotFound,
				Error:   "Task not found",
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

		logger.Error("failed to delete task: ", zap.Error(err))
		return
	}

	resp := response.DeleteResponse{
		Success: true,
		Code:    http.StatusOK,
		Data:    "Task deleted successfully",
	}
	c.JSON(http.StatusOK, resp)
}

func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func derefStatus(s *domain.TaskStatus) domain.TaskStatus {
	if s != nil {
		return *s
	}
	return ""
}
