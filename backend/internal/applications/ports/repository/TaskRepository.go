package repository

import (
	"task-management/internal/domain"
	"time"
)

type TaskRepository interface {
	Create(task *domain.Task) error
	GetByID(id string) (*domain.Task, error)
	GetByUser(userID int, status *domain.TaskStatus, deadline *time.Time) ([]domain.Task, error)
	Update(task *domain.Task) error
	Delete(id string) error
}
