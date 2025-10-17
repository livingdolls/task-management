package storages

import (
	"task-management/internal/applications/ports/repository"
	"task-management/internal/domain"
	"time"

	"gorm.io/gorm"
)

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) repository.TaskRepository {
	return &taskRepository{db: db}
}

// Create implements repository.TaskRepository.
func (t *taskRepository) Create(task *domain.Task) error {
	return t.db.Create(task).Error
}

// Delete implements repository.TaskRepository.
func (t *taskRepository) Delete(id uint) error {
	return t.db.Delete(&domain.Task{}, id).Error
}

// GetByID implements repository.TaskRepository.
func (t *taskRepository) GetByID(id uint) (*domain.Task, error) {
	var task domain.Task

	if err := t.db.First(&task, id).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

// GetByUser implements repository.TaskRepository.
func (t *taskRepository) GetByUser(userID uint, status *domain.TaskStatus, deadline *time.Time) ([]domain.Task, error) {
	var tasks []domain.Task
	query := t.db.Where("user_id = ?", userID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if deadline != nil {
		query = query.Where("deadline <= ?", *deadline).Order("deadline ASC")
	} else {
		query = query.Order("created_at ASC")
	}

	err := query.Find(&tasks).Error
	return tasks, err
}

// Update implements repository.TaskRepository.
func (t *taskRepository) Update(task *domain.Task) error {
	return t.db.Save(task).Error
}
