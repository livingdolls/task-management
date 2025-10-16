package repository

import (
	"task-management/internal/domain"
	"time"
)

type TaskRepository interface {
	Create(task *domain.Task) error
	GetByID(id uint) (*domain.Task, error)
	GetByUser(userID uint, status *domain.TaskStatus, deadline *time.Time) ([]domain.Task, error)
	Update(task *domain.Task) error
	Delete(id uint) error
}
