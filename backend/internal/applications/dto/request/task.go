package request

import (
	"task-management/internal/domain"
	"time"
)

type CreateTask struct {
	Title       string            `json:"title" binding:"required"`
	Description string            `json:"description" binding:"required"`
	Status      domain.TaskStatus `json:"status" binding:"required,oneof='To Do' 'In Progress' 'Done'"`
	Deadline    *time.Time        `json:"deadline,omitempty"`
}

type UpdateTask struct {
	Title       *string            `json:"title,omitempty"`
	Description *string            `json:"description,omitempty"`
	Status      *domain.TaskStatus `json:"status,omitempty" binding:"omitempty"`
	Deadline    *time.Time         `json:"deadline,omitempty"`
}
