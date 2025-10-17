package response

import "time"

type Task struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type BaseTaskResponse struct {
	Success bool `json:"success"`
	Code    int  `json:"code"`
	Data    Task `json:"data"`
}

type ListTaskResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Data    []Task `json:"data"`
}

type DeleteResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Data    string `json:"data"`
}
