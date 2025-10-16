package services

import (
	"errors"
	"task-management/internal/applications/ports/repository"
	"task-management/internal/applications/ports/services"
	"task-management/internal/domain"
	"time"

	"gorm.io/gorm"
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
	req.UserID = userId
	req.CreatedBy = userId

	if req.Status == "" {
		req.Status = domain.ToDo
	}

	return t.taskRepo.Create(req)
}

// DeleteTask implements services.TaskService.
func (t *taskService) DeleteTask(taskId uint, userId uint) error {
	task, err := t.taskRepo.GetByID(taskId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("task not found")
		}

		return err
	}

	if task.UserID != userId {
		return errors.New("unauthorized")
	}

	return t.taskRepo.Delete(taskId)
}

// GetTasks implements services.TaskService.
func (t *taskService) GetTasks(userId uint, status *domain.TaskStatus, deadline *time.Time) ([]domain.Task, error) {
	return t.taskRepo.GetByUser(userId, status, deadline)
}

// UpdateTask implements services.TaskService.
func (t *taskService) UpdateTask(arg *domain.Task, userId uint) error {
	taskInDb, err := t.taskRepo.GetByID(arg.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("task not found")
		}
		return err
	}

	if taskInDb.UserID != userId {
		return errors.New("unauthorized")
	}

	taskInDb.Title = arg.Title
	taskInDb.Description = arg.Description
	taskInDb.Status = arg.Status
	taskInDb.Deadline = arg.Deadline

	return t.taskRepo.Update(taskInDb)
}

// GetTaskById implements services.TaskService.
func (t *taskService) GetTaskById(taskId uint, userId uint) (*domain.Task, error) {
	task, err := t.taskRepo.GetByID(taskId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	if task.UserID != userId {
		return nil, errors.New("unauthorized")
	}

	return task, nil
}
