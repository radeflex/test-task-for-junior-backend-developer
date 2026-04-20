package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type taskMutationDTO struct {
	Title            string                  `json:"title"`
	Description      string                  `json:"description"`
	Status           taskdomain.Status       `json:"status"`
	Periodicity      *taskdomain.Periodicity `json:"periodicity,omitempty"`
	PeriodicityValue *int                    `json:"periodicity_value,omitempty"`
	PublishDate      *time.Time              `json:"publish_date,omitempty"`
}

type taskDTO struct {
	ID               int64                   `json:"id"`
	Title            string                  `json:"title"`
	Description      string                  `json:"description"`
	Status           taskdomain.Status       `json:"status"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
	Periodicity      *taskdomain.Periodicity `json:"periodicity,omitempty"`
	PeriodicityValue *int                    `json:"periodicity_value,omitempty"`
	PublishDate      *time.Time              `json:"publish_date,omitempty"`
	IsActive         bool                    `json:"is_active"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:               task.ID,
		Title:            task.Title,
		Description:      task.Description,
		Status:           task.Status,
		CreatedAt:        task.CreatedAt,
		UpdatedAt:        task.UpdatedAt,
		Periodicity:      task.Periodicity,
		PeriodicityValue: task.PeriodicityValue,
		PublishDate:      task.PublishDate,
		IsActive:         task.IsActive,
	}
}
