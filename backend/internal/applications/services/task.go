package services

import (
	"task-management/internal/applications/ports/repository"
	"task-management/internal/applications/ports/services"
	"task-management/internal/domain"
	"time"
)

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) services.TaskService {
	return &taskService{
		taskRepo: repo,
	}
}

// CreateTask implements services.TaskService.
func (t *taskService) CreateTask(userId uint, req *domain.Task) error {
	panic("unimplemented")
}

// DeleteTask implements services.TaskService.
func (t *taskService) DeleteTask(taskId uint) error {
	panic("unimplemented")
}

// GetTasks implements services.TaskService.
func (t *taskService) GetTasks(userId uint, status *domain.TaskStatus, deadline *time.Time) ([]domain.Task, error) {
	panic("unimplemented")
}

// UpdateTask implements services.TaskService.
func (t *taskService) UpdateTask(task *domain.Task) error {
	panic("unimplemented")
}
