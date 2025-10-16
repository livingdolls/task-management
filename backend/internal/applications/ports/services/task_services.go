package services

import (
	"task-management/internal/domain"
	"time"
)

type TaskService interface {
	CreateTask(userId uint, req *domain.Task) error
	GetTasks(userId uint, status *domain.TaskStatus, deadline *time.Time) ([]domain.Task, error)
	UpdateTask(task *domain.Task) error
	DeleteTask(taskId uint) error
}
