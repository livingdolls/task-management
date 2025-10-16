package domain

import "time"

type TaskStatus string

const (
	ToDo       TaskStatus = "To Do"
	InProgress TaskStatus = "In Progress"
	Done       TaskStatus = "Done"
)

type Task struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `gorm:"index;not null" json:"user_id"`
	Title       string     `gorm:"size:255;not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Status      TaskStatus `gorm:"size:20;not null;default:'To Do'" json:"status"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	CreatedBy   uint       `json:"created_by"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
}
