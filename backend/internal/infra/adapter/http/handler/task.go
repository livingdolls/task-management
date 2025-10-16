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

func (h *TaskHandler) Create(c *gin.Context) {
	var req request.CreateTask

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"error":   err.Error(),
		})
		return
	}

	// claims token dari middleware
	userClaims, ok := middleware.GetUserClaims(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"error":   "Unauthorized",
		})
		return
	}

	task := domain.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Deadline:    req.Deadline,
	}

	if err := h.taskService.CreateTask(userClaims.UserID, &task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"error":   "Internal server error",
		})

		logger.Info("failed to create task: ", zap.Error(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"code":    http.StatusCreated,
		"data": response.Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      string(task.Status),
			Deadline:    task.Deadline,
			CreatedAt:   task.CreatedAt,
		},
	})
}

func (h *TaskHandler) Get(c *gin.Context) {
	// claims token dari middleware
	userClaims, ok := middleware.GetUserClaims(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"error":   "Unauthorized",
		})
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
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"code":    http.StatusBadRequest,
				"error":   "Invalid deadline format. Use YYYY-MM-DD",
			})
			return
		}
		deadline = &parsedDeadline
	}

	tasks, err := h.taskService.GetTasks(userClaims.UserID, status, deadline)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"error":   "Internal server error",
		})

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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"data":    resp,
	})
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"error":   "Invalid task ID",
		})
		return
	}

	// claims token dari middleware
	userClaims, ok := middleware.GetUserClaims(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"error":   "Unauthorized",
		})
		return
	}

	task, err := h.taskService.GetTaskById(uint(id), userClaims.UserID)

	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"error":   "Unauthorized",
			})
			return
		}

		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"code":    http.StatusNotFound,
				"error":   "Task not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"error":   "Internal server error",
		})

		logger.Info("failed to get task by id: ", zap.Error(err))
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"code":    http.StatusNotFound,
			"error":   "Task not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"data": response.Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      string(task.Status),
			Deadline:    task.Deadline,
			CreatedAt:   task.CreatedAt,
		},
	})
}

func (h *TaskHandler) Update(c *gin.Context) {
	var req request.UpdateTask

	// Get task ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"error":   "Invalid task ID",
		})
		return
	}

	// claims token dari middleware
	userClaims, ok := middleware.GetUserClaims(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"error":   "Unauthorized",
		})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"error":   err.Error(),
		})
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
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"error":   "Unauthorized",
			})
			return
		}

		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"code":    http.StatusNotFound,
				"error":   "Task not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"error":   "Internal server error",
		})

		logger.Error("failed to update task: ", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"data": response.Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      string(task.Status),
			Deadline:    task.Deadline,
			CreatedAt:   task.CreatedAt,
		},
	})
}

func (h *TaskHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    http.StatusBadRequest,
			"error":   "Invalid task ID",
		})
		return
	}

	// claims token dari middleware
	userClaims, ok := middleware.GetUserClaims(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"error":   "Unauthorized",
		})
		return
	}

	if err := h.taskService.DeleteTask(uint(id), userClaims.UserID); err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"error":   "Unauthorized",
			})
			return
		}

		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"code":    http.StatusNotFound,
				"error":   "Task not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    http.StatusInternalServerError,
			"error":   "Internal server error",
		})

		logger.Error("failed to delete task: ", zap.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"data":    "Task deleted successfully",
	})
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
