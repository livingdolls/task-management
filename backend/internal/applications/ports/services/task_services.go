package services

import (
	"task-management/internal/domain"
	"time"
)

type TaskService interface {
	CreateTask(userId uint, req *domain.Task) error
	GetTasks(userId uint, status *domain.TaskStatus, deadline *time.Time) ([]domain.Task, error)
	UpdateTask(arg *domain.Task, userId uint) error
	DeleteTask(taskId uint, userId uint) error
}
